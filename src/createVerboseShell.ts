import { $, ShellPromise, type ShellExpression } from "bun";

export function createVerboseShell(
	dryRun: boolean | undefined,
	verbose: boolean | undefined
): (strings: TemplateStringsArray, ...expressions: ShellExpression[]) => ShellPromise | void {
	return function $debug(strings, ...expressions) {
		if (dryRun || verbose) {
			let result = "";
			for (let i = 0; i < strings.length; i++) {
				result += strings[i];
				if (i < expressions.length) {
					result += $.escape(String(expressions[i]));
				}
			}
			console.log(result);
		}
		if (!dryRun) {
			return $(strings, ...expressions);
		}
	};
}
