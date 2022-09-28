// use std::path::Path;
// use std::fs;

// fn walk(path: &str) {
//     // check if path is folder or file
//     if Path::new(path).is_dir() {
//         // if folder, iterate over all files in folder
//         let result = fs::read_dir(path);
//         match result {
//             Ok(paths) => {
//                 for path in paths {
//                     // call walk function on each file
//                     walk(&path.unwrap().path().to_str().unwrap());
//                 }
//             },
//             Err(e) => println!("Error: {}", e),
//         }
//     } else {
//         check_path(path);
//     }
// }

fn check_path(path: &str) {
    let python = format!("{}/python", path);
    let python3 = format!("{}/python3", path);
    if std::path::Path::new(&python).exists() {
        println!("Found python at {}", python);
    }
    if std::path::Path::new(&python3).exists() {
        println!("Found python3 at {}", python3);
    }
}

fn main() {
    // array of common executable folders on Unix
    let known_paths = [
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
    ];

    let path = std::env::var("PATH").unwrap();
    let paths_from_env = path.split(":").chain(known_paths.iter().map(|x| *x));

    let mut handles = vec![];
    for path in paths_from_env {
        let path = path.to_string();
        let handle = std::thread::spawn(move || {
            check_path(&path);
        });
        handles.push(handle);
    }

    // search for python or python3 in .pyenv versions folder
    let home = std::env::var("HOME").unwrap();
    let pyenv = format!("{}/.pyenv", home);
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

    for handle in handles {
        handle.join().unwrap();
    }
}
