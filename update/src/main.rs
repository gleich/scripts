use std::{env, path::Path, process};

use anyhow::{Context, Result};
use task_log::{task, ConfigBuilder};

fn main() -> Result<()> {
    ConfigBuilder::new()
        .replace(false)
        .apply()
        .expect("Failed to apply custom logging configuration");
    let args: Vec<String> = env::args().collect();

    command("brew", vec!["update"], None)?;
    command("brew", vec!["upgrade"], None)?;
    command("brew", vec!["cleanup", "-s"], None)?;
    command("rustup", vec!["update"], None)?;
    // command("cargo", vec!["install-update", "-a"], None)?;

    let rmapi_dir = Path::new("/Users/matt/tmp/rmapi");
    command("git", vec!["pull"], Some(rmapi_dir))?;
    command("go", vec!["install", "."], Some(rmapi_dir))?;

    if args.contains(&String::from("--fetch")) {
        command("fetch", vec![], Some(Path::new("/Users/matt/src/dots")))?;
    }

    Ok(())
}

fn command(binary: &str, args: Vec<&str>, path: Option<&Path>) -> Result<()> {
    let cmd = format!("{} {}", binary, args.join(" "));
    task(&cmd, || -> Result<()> {
        let mut process = process::Command::new(binary)
            .args(&args)
            .current_dir(path.unwrap_or(Path::new("/Users/matt")))
            .spawn()
            .context(format!("Failed to run command: {}", cmd))?;
        process
            .wait()
            .context(format!("Failed to wait for {} to finish", cmd))?;
        Ok(())
    })?;
    Ok(())
}
