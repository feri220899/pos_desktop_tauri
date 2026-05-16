use std::{collections::HashMap, fs, path::PathBuf, sync::Mutex};

use gethostname::gethostname;
use hex::encode as hex_encode;
use jsonwebtoken::{decode, Algorithm, DecodingKey, Validation};
use mac_address::get_mac_address;
use serde::{Deserialize, Serialize};
use serde_json::{json, Value};
use sha2::{Digest, Sha256};
use tauri::{Manager, State};
use tauri_plugin_shell::ShellExt;

struct SidecarState(Mutex<Option<tauri_plugin_shell::process::CommandChild>>);

const RSA_PUBLIC_KEY: &str = "-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvGDG8GdB6J63/23CEFuV
KL6Y0GFpBx1DC6RyVmyaY8jiJCRymKOk0QjwoGuubQw2izu2zb6lc5uj4XibrSB6
l2QDrcx7gZ3Hj+QOUBtCwTrCO8FaGCHHJH28qtwvqCzsXvLxC6qacEMEikLt11EW
pK8PlnvYPwZrSCf/GtZfH3a+T/Z2YmbAgePLXUzlmAmzWAK8jvLjQaqQMCPAMUr9
mHnzrKyVqh9Gi8Iy41RUWlUGlrfJoY0oz6X3QBzr6jspuRS/dg2KIPTqlpAPwFt1
CvqY3QJ4nn6j7s7GwT/mQJ6s4oKX1+OW679O10UkyUP2uebaYhTRS4dgQE79u0uC
ZQIDAQAB
-----END PUBLIC KEY-----";

#[derive(Debug, Serialize, Deserialize)]
struct LicenseClaims {
    device_id: Option<String>,
    license_key: Option<String>,
    tipe: Option<String>,
    exp: Option<u64>,
    #[serde(flatten)]
    extra: HashMap<String, Value>,
}

fn config_path(app: &tauri::AppHandle) -> PathBuf {
    app.path()
        .app_data_dir()
        .expect("failed to get app data dir")
        .join("config.json")
}

fn read_config(app: &tauri::AppHandle) -> Value {
    let path = config_path(app);
    if let Ok(content) = fs::read_to_string(&path) {
        serde_json::from_str(&content).unwrap_or(json!({}))
    } else {
        json!({})
    }
}

fn write_config(app: &tauri::AppHandle, config: &Value) -> Result<(), String> {
    let path = config_path(app);
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent).map_err(|e| e.to_string())?;
    }
    let content = serde_json::to_string_pretty(config).map_err(|e| e.to_string())?;
    fs::write(&path, content).map_err(|e| e.to_string())
}

#[tauri::command]
fn get_config(app: tauri::AppHandle, key: String) -> Value {
    let config = read_config(&app);
    config.get(&key).cloned().unwrap_or(Value::Null)
}

#[tauri::command]
fn set_config(app: tauri::AppHandle, key: String, value: Value) -> Result<(), String> {
    let mut config = read_config(&app);
    config[key] = value;
    write_config(&app, &config)
}

#[tauri::command]
fn get_device_id(app: tauri::AppHandle) -> Result<String, String> {
    let config = read_config(&app);
    if let Some(id) = config.get("device_id").and_then(|v| v.as_str()) {
        return Ok(id.to_string());
    }

    let mac = get_mac_address()
        .map_err(|e| e.to_string())?
        .map(|m| m.to_string())
        .unwrap_or_else(|| "00:00:00:00:00:00".to_string());

    let hostname = gethostname().to_string_lossy().to_string();
    let raw = format!("{}-{}", mac, hostname);

    let mut hasher = Sha256::new();
    hasher.update(raw.as_bytes());
    let device_id = hex_encode(hasher.finalize());

    let mut config = read_config(&app);
    config["device_id"] = json!(device_id);
    write_config(&app, &config)?;

    Ok(device_id)
}

#[tauri::command]
fn verify_license_token(app: tauri::AppHandle, token: String) -> Value {
    let decoding_key = match DecodingKey::from_rsa_pem(RSA_PUBLIC_KEY.as_bytes()) {
        Ok(k) => k,
        Err(e) => {
            return json!({ "valid": false, "expired": false, "error": e.to_string() });
        }
    };

    let mut validation = Validation::new(Algorithm::RS256);
    validation.validate_exp = true;

    match decode::<LicenseClaims>(&token, &decoding_key, &validation) {
        Ok(token_data) => {
            let claims = token_data.claims;

            let stored_device_id = get_device_id(app)
                .unwrap_or_default();

            let token_device_id = claims.device_id.as_deref().unwrap_or("");
            if token_device_id != stored_device_id {
                return json!({
                    "valid": false,
                    "expired": false,
                    "error": "device_id mismatch"
                });
            }

            let now = std::time::SystemTime::now()
                .duration_since(std::time::UNIX_EPOCH)
                .map(|d| d.as_secs())
                .unwrap_or(0);

            let days_left = claims.exp.map(|exp| {
                if exp > now {
                    ((exp - now) as f64 / 86400.0).ceil() as i64
                } else {
                    0
                }
            });

            json!({
                "valid": true,
                "expired": false,
                "days_left": days_left,
                "payload": {
                    "device_id": claims.device_id,
                    "license_key": claims.license_key,
                    "tipe": claims.tipe,
                    "exp": claims.exp,
                }
            })
        }
        Err(e) => {
            let is_expired = e.to_string().contains("ExpiredSignature");

            if is_expired {
                json!({ "valid": false, "expired": true, "error": "token expired" })
            } else {
                json!({ "valid": false, "expired": false, "error": e.to_string() })
            }
        }
    }
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_shell::init())
        .manage(SidecarState(Mutex::new(None)))
        .setup(|app| {
            let handle   = app.handle();
            let data_dir = handle.path().app_data_dir().expect("no app data dir");
            let db_path  = data_dir.join("data.db").to_string_lossy().to_string();

            let config    = read_config(&handle);
            let app_mode  = config.get("app_mode").and_then(|v| v.as_str()).unwrap_or("").to_string();

            let (_rx, child) = handle
                .shell()
                .sidecar("pos-backend")
                .expect("pos-backend sidecar not configured")
                .args(["--port", "3001", "--db", &db_path, "--mode", &app_mode])
                .spawn()
                .expect("failed to spawn pos-backend");

            let state: State<SidecarState> = handle.state();
            *state.0.lock().unwrap() = Some(child);
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            get_config,
            set_config,
            get_device_id,
            verify_license_token,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
