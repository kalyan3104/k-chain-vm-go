package vmhooksgenerate

import (
	"fmt"
)

// WriteRustVMHooksTrait autogenerate data in the provided file
func WriteRustVMHooksTrait(out *eiGenWriter, eiMetadata *EIMetadata) {
	autoGeneratedHeader(out)
	out.WriteString(`
use std::ffi::c_void;

use crate::{MemLength, MemPtr};

#[rustfmt::skip]
#[allow(clippy::too_many_arguments)]
pub trait VMHooks: core::fmt::Debug + 'static {
    fn set_vm_hooks_ptr(&mut self, vm_hooks_ptr: *mut c_void);

`)

	for _, funcMetadata := range eiMetadata.AllFunctions {
		out.WriteString(fmt.Sprintf(
			"    fn %s%s;\n",
			snakeCase(funcMetadata.Name),
			writeRustFnDeclarationArguments(
				"&self",
				funcMetadata,
				rustVMHooksType,
			),
		))
	}

	out.WriteString(`}

/// Dummy implementation for VMHooks. Can be used as placeholder, or in tests.
#[derive(Debug)]
pub struct VMHooksDefault;

#[allow(unused)]
#[rustfmt::skip]
impl VMHooks for VMHooksDefault {
    fn set_vm_hooks_ptr(&mut self, _vm_hooks_ptr: *mut c_void) {
    }

`)

	for i, funcMetadata := range eiMetadata.AllFunctions {
		if i > 0 {
			out.WriteString("\n")
		}

		out.WriteString(fmt.Sprintf(
			"    fn %s%s {\n",
			snakeCase(funcMetadata.Name),
			writeRustFnDeclarationArguments(
				"&self",
				funcMetadata,
				rustVMHooksType,
			),
		))

		out.WriteString(fmt.Sprintf(
			"        println!(\"Called: %s\");\n",
			snakeCase(funcMetadata.Name),
		))

		if funcMetadata.Result != nil {
			out.WriteString("        0\n")
		}

		out.WriteString("    }\n")
	}

	out.WriteString("}\n")

}