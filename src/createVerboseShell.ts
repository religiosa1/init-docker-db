import { $, ShellPromise, type ShellExpression } from "bun";

export type VerboseShell = (
	strings: TemplateStringsArray,
	...expressions: ShellExpression[]
) => ShellPromise | undefined;

export function createVerboseShell(dryRun: boolean | undefined, verbose: boolean | undefined): VerboseShell {
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
