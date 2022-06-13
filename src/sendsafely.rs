use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct Confirmation {
    ip_address: String,
    timestamp: String,
    time_stamp_str: String,
    is_message: bool,
}

#[derive(Deserialize, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct Recipient {
    recipient_id: String,
    email: String,
    full_name: String,
    needs_approval: bool,
    recipient_code: String,
    confirmations: Vec<Confirmation>,
    is_package_owner: bool,
    check_for_public_keys: bool,
    role_name: String,
}

#[derive(Deserialize, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ObjectWithId {
    id: String,
}

#[derive(Deserialize, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct Package {
    package_id: String,
    package_code: String,
    server_secret: String,
    recipients: Vec<Recipient>,
    contact_groups: Vec<ObjectWithId>,
    files: Vec<ObjectWithId>,
    directories: Vec<ObjectWithId>,
    approver_list: Vec<ObjectWithId>,
    needs_approval: bool,
    state: String,
    password_required: bool,
    life: i128,
    #[serde(rename = "isVDR")]
    is_vdr: bool,
    is_archived: bool,
    package_sender: String,
    package_timestamp: String,
    root_directory_id: String,
    response: String,
}
