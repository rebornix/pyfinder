use std::fs::File;
use std::io::{self, BufRead};
use std::path::Path;
fn check_path(path: &str) {
    if Path::new(path).is_dir() {
        // println!("Checking {} ...", path);

        for entry in std::fs::read_dir(path).unwrap() {
            let entry = entry.unwrap();
            let path = entry.path();
            let path = path.to_str().unwrap();
            // check if file is python or python3
            // if it's windows
            if cfg!(windows) {
                if path.ends_with("python.exe") || path.ends_with("python3.exe") {
                    println!("{}", path);
                }
            } else {
                if path.ends_with("python") || path.ends_with("python3") {
                    println!("{}", path);
                }
            }
            
        }
    }
}

fn main() {
    // array of common executable folders on Unix
    // paths based on os
    let known_paths = if cfg!(windows) {
        vec![
        ]
    } else {
        vec![
            "/bin",
            "/etc",
            "/lib",
            "/lib/x86_64-linux-gnu",
            "/lib64",
            "/sbin",
            "/snap/bin",
            "/usr/bin",
            "/usr/games",
            "/usr/include",
            "/usr/lib",
            "/usr/lib/x86_64-linux-gnu",
            "/usr/lib64",
            "/usr/libexec",
            "/usr/local",
            "/usr/local/bin",
            "/usr/local/etc",
            "/usr/local/games",
            "/usr/local/lib",
            "/usr/local/sbin",
            "/usr/sbin",
            "/usr/share",
            "~/.local/bin"
        ]
    };

    let path = std::env::var("PATH").unwrap();
    // path list separator
    let separator = if cfg!(windows) { ";" } else { ":" };
    let paths_from_env = path.split(separator).chain(known_paths.iter().map(|x| *x));
    // println!(paths_from_env)
    let mut handles = vec![];
    for path in paths_from_env {
        let path = path.to_string();
        // println!("Checking {} ...", path);
        let handle = std::thread::spawn(move || {
            check_path(&path);
        });
        handles.push(handle);
    }

    // search for python or python3 in .pyenv versions folder
    let home_env = std::env::var("HOME");
    if home_env.is_ok() {
        let pyenv = format!("{}/.pyenv", home_env.unwrap());
        let versions = format!("{}/versions", pyenv);
        for entry in std::fs::read_dir(versions).unwrap() {
            let entry = entry.unwrap();
            let path = entry.path();
            let bin = format!("{}/bin", path.display());
            let handle = std::thread::spawn(move || {
                check_path(&bin);
            });
            handles.push(handle);
        }
    }

    // check if it's Windows
    if cfg!(windows) {
        // check localAppData
        let local_app_data = std::env::var("LOCALAPPDATA");
        // print
        if local_app_data.is_ok() {
            // windows store root
            let windows_store_root = format!("{}\\Microsoft\\WindowsApps", local_app_data.unwrap());
            let handle = std::thread::spawn(move || {
                check_path(&windows_store_root);
            });
            handles.push(handle);
        }
    }

    // check for conda
    // let conda_environments_path = format(args)
    // concat path fragments
    let conda_environments_path = if cfg!(windows) {
        format!("{}\\.conda\\environments.txt", std::env::var("USERPROFILE").unwrap())
    } else {
        format!("{}/.conda/environments.txt", std::env::var("HOME").unwrap())
    };

    // check if file exists
    if Path::new(&conda_environments_path).exists() {
        // read file
        if let Ok(lines) = read_lines(conda_environments_path) {
            // Consumes the iterator, returns an (Optional) String
            for line in lines {
                if let Ok(ip) = line {
                    let handle = std::thread::spawn(move || {
                        check_path(&ip);
                    });
                    handles.push(handle);
                }
            }
        }
    }

    for handle in handles {
        handle.join().unwrap();
    }
}
fn read_lines<P>(filename: P) -> io::Result<io::Lines<io::BufReader<File>>>
where P: AsRef<Path>, {
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}
