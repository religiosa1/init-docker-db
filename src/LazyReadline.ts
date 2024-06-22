import readline from "node:readline/promises";

export class LazyReadline implements Disposable {
	private rl: readline.Interface | undefined;

	constructor(private disabled: boolean) {}

	question(query: string): Promise<string> {
		return this.get().question(query);
	}

	[Symbol.dispose](): void {
		this.rl?.close();
	}

	private get(): readline.Interface {
		if (this.rl) {
			return this.rl;
		}
		if (this.disabled) {
			throw new ReadlineDisabledError();
		}
		this.rl = readline.createInterface({
			input: process.stdin,
			output: process.stdout,
		});
		return this.rl;
	}
}

export class ReadlineDisabledError extends Error {}
