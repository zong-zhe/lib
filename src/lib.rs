//! # KCL Rust SDK
//!
//! ## How to Use
//!
//! ```no_check,no_run
//! cargo add --git https://github.com/kcl-lang/lib --branch main
//! ```
//!
//! Write the Code
//!
//! ```rust
//! use kcl_lang::*;
//! use std::path::Path;
//! use anyhow::Result;
//!
//! fn main() -> Result<()> {
//!     let api = API::default();
//!     let args = &ExecProgramArgs {
//!         work_dir: Path::new(".").join("testdata").canonicalize().unwrap().display().to_string(),
//!         k_filename_list: vec!["test.k".to_string()],
//!         ..Default::default()
//!     };
//!     let exec_result = api.exec_program(args)?;
//!     assert_eq!(exec_result.yaml_result, "alice:\n  age: 18");
//!     Ok(())
//! }
//! ```

use std::ffi::{CStr, CString};

pub use kclvm_api::gpyrpc::*;
use kclvm_api::service::capi::{kclvm_service_call_with_length, kclvm_service_new};
use kclvm_api::service::service_impl::KclvmServiceImpl;

use anyhow::Result;

pub type API = KclvmServiceImpl;

pub fn call<'a>(name: &'a [u8], args: &'a [u8]) -> Result<&'a [u8]> {
    let mut result_len: usize = 0;
    let result_ptr = {
        let args = CString::new(args)?;
        let call = CString::new(name)?;
        let serv = kclvm_service_new(0);
        kclvm_service_call_with_length(serv, call.as_ptr(), args.as_ptr(), &mut result_len)
    };
    // let result = unsafe { CStr::from_ptr(result_ptr) };
    let result_slice = unsafe { std::slice::from_raw_parts(result_ptr as *const u8, result_len + 1) }; // 包含 null 字符
    let result_cstr = CStr::from_bytes_with_nul(result_slice)?;
    Ok(result_cstr.to_bytes())
}
