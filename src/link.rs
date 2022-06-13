use time::{macros::format_description, OffsetDateTime};

use curl::easy::{Easy, List};

use crate::config::SendSafelyConfig;

const BASE_URL: &str = "https://demo.sendsafely.com/api/v2.0";

pub fn extract_package(url: String) -> String {
    return "".to_string();
}
pub fn download_link(url: String, config: SendSafelyConfig) {
    let package_id = extract_package(url);
    let get_packages_url = [BASE_URL.to_string(), "package".to_string(), package_id].join("/");

    let mut easy = Easy::new();
    easy.url(&get_packages_url).unwrap();

    let mut list = List::new();
    list.append("Content-Type: application/json").unwrap();
    let ss_api_key = format!("ss-api-key: {}", config.ss_api_key);
    list.append(ss_api_key.as_str()).unwrap();
    let req_ts = OffsetDateTime::now_utc();
    let fmt =
        format_description!("[month repr:short] [day], [year] [hour]:[minute]:[second] [period]");
    let formatted = req_ts.format(fmt);
    let ss_api_key = format!("ss-request-timestamp: {}", formatted.unwrap());
    list.append(ss_api_key.as_str()).unwrap();
    easy.http_headers(list).unwrap();
    easy.perform().unwrap();
}
