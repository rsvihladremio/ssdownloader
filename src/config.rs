use std::{fs, path::Path};

use serde::{Deserialize, Serialize};
use toml::de::Error as TomlError;

#[derive(Deserialize, Serialize)]
pub struct SendSafelyConfig {
    pub ss_api_key: String,
    pub ss_api_secret: String,
}

#[derive(Deserialize, Serialize)]
pub struct ZendeskConfig {
    pub zendesk_url: String,
}

pub fn save_config(ss_api_key: &String, ss_api_secret: &String, zendesk_url: &String) {
    std::fs::create_dir_all(Path::new(config_dir().as_str())).expect("unable to create directory");
    let ss_config = SendSafelyConfig {
        ss_api_key: ss_api_key.to_owned(),
        ss_api_secret: ss_api_secret.to_owned(),
    };
    let ss_text =
        toml::to_string(&ss_config).expect("unable to serialize send safely configuration file");
    fs::write(config_file_ss(), ss_text).expect("error writing send safely configuration file");
    let zd_config = ZendeskConfig {
        zendesk_url: zendesk_url.to_owned(),
    };
    let zd_text =
        toml::to_string(&zd_config).expect("unable to serialize zendesk configuration file");
    fs::write(config_file_zendesk(), zd_text).expect("error writing zendesk configuration file");
}

fn config_file_ss() -> String {
    [config_dir(), "ssendsafely.toml".to_string()].join("/")
}

fn config_dir() -> String {
    let home_dir = dirs::home_dir().expect("unable to read home dir");
    let home_dir_str = home_dir
        .to_str()
        .expect("unable convert home dir path to string");
    [
        home_dir_str.to_string(),
        ".config".to_string(),
        "ssdownloader".to_string(),
    ]
    .join("/")
}

fn config_file_zendesk() -> String {
    [config_dir(), "zendesk.toml".to_string()].join("/")
}

pub fn load_ss() -> SendSafelyConfig {
    let config_ss = config_file_ss();
    let contents = fs::read_to_string(config_ss)
        .expect("Something went wrong reading the send safely configuration file");

    let ss_config_result: Result<SendSafelyConfig, TomlError> = toml::from_str(&contents);
    ss_config_result.expect("unable to deserialize send safely configuration")
}

pub fn load_zd() -> ZendeskConfig {
    let config_zd = config_file_zendesk();
    let contents = fs::read_to_string(config_zd)
        .expect("Something went wrong reading the zendesk configuration file");

    let zd_config_result: Result<ZendeskConfig, TomlError> = toml::from_str(&contents);
    zd_config_result.expect("unable to deserialize zendesk configuration")
}
