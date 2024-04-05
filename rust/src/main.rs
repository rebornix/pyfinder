use std::fs::File;
use std::io::{self, BufRead};
use std::path::Path;

const PYTHON_SUFFIX_ARRAY: [&str; 4] = if cfg!(windows) {
    ["python.exe", "python3.exe", "ipython.exe", "ipython3.exe"]
} else {
    ["python", "python3", "ipython", "ipython3"]
};

fn check_path(path: &Path) {
    for suffix in PYTHON_SUFFIX_ARRAY.iter() {
        let python_path = path.join(suffix);
        if python_path.exists() {
            println!("{}", python_path.display());
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
        let path = Path::new(path).to_path_buf();
        // println!("Checking {} ...", path);
        let handle = std::thread::spawn(move || {
            check_path(&path);
        });
        handles.push(handle);
    }

    // search for python or python3 in .pyenv versions folder
    let home_env = std::env::var("HOME");
    if home_env.is_ok() {
        let pyenv = Path::new(&home_env.unwrap()).join(".pyenv");
        let versions = pyenv.join("versions");
        for entry in std::fs::read_dir(versions).unwrap() {
            let entry = entry.unwrap();
            let path = entry.path();
            let bin = path.join("bin");
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
            let windows_store_root = Path::new(&local_app_data.unwrap()).join("Microsoft").join("WindowsApps");
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
                    let ip_path = Path::new(&ip).to_path_buf(); // Clone the value of `ip`
                    let handle = std::thread::spawn(move || {
                        check_path(&ip_path);
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
