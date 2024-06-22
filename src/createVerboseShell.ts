import { $, ShellPromise, type ShellExpression } from "bun";

export function createVerboseShell(
	verbose?: boolean
): (strings: TemplateStringsArray, ...expressions: ShellExpression[]) => ShellPromise {
	return function $debug(strings, ...expressions) {
		if (verbose) {
			let result = "";
			for (let i = 0; i < strings.length; i++) {
				result += strings[i];
				if (i < expressions.length) {
					result += $.escape(String(expressions[i]));
				}
			}
			console.log("CMD:", result);
		}
		return $(strings, ...expressions);
	};
}
