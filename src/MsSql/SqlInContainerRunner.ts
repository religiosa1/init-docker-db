import type { ShellOutput } from "bun";
import type { VerboseShell } from "../createVerboseShell";
import type { IDbCreateOptions } from "../DbCreator";
import { escapeId } from "./escape";

export class SqlInContainerRunner {
	constructor(
		private readonly contId: string,
		private readonly shell: VerboseShell,
		private readonly opts: IDbCreateOptions
	) {}

	async run(sql: string): Promise<boolean> {
		if (this.opts.verbose) {
			console.log("SQL:", sql.includes("\n") ? "\n" + sql + " --> END SQL" : sql);
		}
		let prms = this.shell`docker exec ${this.contId} \
/opt/mssql-tools/bin/sqlcmd -S localhost \
-U SA -P ${this.opts.password} -Q ${sql} || exit 1`;
		if (!this.opts.verbose && !this.opts.dryRun) {
			prms = prms?.quiet();
		}

		let output: ShellOutput | undefined;
		try {
			output = await prms;
		} catch {
			return false;
		}

		if (this.opts.dryRun || !output) return true;
		const outputText = output.text();
		if (this.#checkOutputForSqlErrors(outputText)) {
			if (!this.opts.verbose) {
				console.log(outputText);
			}
			throw new SqlCommandError(sql, outputText);
		}
		return true;
	}

	runInDb(sql: string): Promise<boolean> {
		return this.run(`use ${escapeId(this.opts.database)}\n` + sql);
	}

	#checkOutputForSqlErrors(output: string): boolean {
		const match = output.match(/^Msg (?:\d+), Level (\d+), State (?:\d+), Server (?:[^,]+), Line (?:\d+)/);
		if (match == null) return false;
		const level = +match[1];
		// Any severity level less or equal 10 we treat as not an error
		// https://learn.microsoft.com/en-us/sql/relational-databases/errors-events/database-engine-error-severities?view=sql-server-ver16#levels-of-severity
		return level > 10;
	}
}

class SqlCommandError extends Error {
	override name = "SqlCommandError";

	constructor(private readonly sql: string, private readonly databaseEngineOutput: string) {
		super("Sql command error");
	}

	override toString() {
		return this.message + ":\n" + this.sql + "\n" + this.databaseEngineOutput;
	}
}
