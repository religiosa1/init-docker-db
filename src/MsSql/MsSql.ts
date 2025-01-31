import { DbCreator, type PasswordValidityTuple } from "../DbCreator";

import { MsSqlPwdValidityEnum, msSqlPwdValidityEnumMessage } from "./MsSqlPwdValidityEnum";
import { escapeId, escapeUser, escapeStr } from "./escape";
import { waitFor } from "./waitFor";

// https://learn.microsoft.com/en-us/sql/relational-databases/security/password-policy?view=sql-server-ver16#password-complexity
function validatePasswordComplexity(password: string): boolean {
	let numberOfCharClassesMatched = 0;

	if (/[A-Z]/.test(password)) {
		numberOfCharClassesMatched++;
	}
	if (/[a-z]/.test(password)) {
		numberOfCharClassesMatched++;
	}
	if (/\d/.test(password)) {
		numberOfCharClassesMatched++;
	}
	if (/[!@#$%^&*()_\-+={}\[\]\\|/<>~,.;:'"]/.test(password)) {
		numberOfCharClassesMatched++;
	}

	return numberOfCharClassesMatched >= 3;
}

export const MsSql = new DbCreator({
	name: "mssql",
	port: 1433,
	defaultUser: "mssql",
	defaultTag: "2022-latest",
	defaultPassword: "Password12",

	isPasswordValid(password: string): PasswordValidityTuple {
		if (!password.length) {
			return [false, msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.Empty]];
		}
		if (password.length < 10) {
			return [false, msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.TooShort]];
		}
		if (!validatePasswordComplexity(password)) {
			return [false, msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.TooSimple]];
		}
		return [true, msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.Valid]];
	},

	async create($, opts) {
		const vlog = (...args: unknown[]) => (opts.verbose ? console.log("verbose:", ...args) : () => {});

		// https://mcr.microsoft.com/product/mssql/server/about
		const shellOutput = await $`docker run -e ACCEPT_EULA=Y \
--name ${opts.containerName}\
--hostname ${opts.containerName}\
-e MSSQL_SA_PASSWORD=${opts.password}\
-p ${this.port}:${opts.port}\
-d mcr.microsoft.com/mssql/server:${opts.tag}`;
		const contId = shellOutput?.text().trim();

		// TODO: parse the resposnse os SQL server and check for potential errors
		const sqlcmd = async (sql: string, db?: string) => {
			if (db) {
				sql = `use ${escapeId(opts.database)}\n` + sql;
			}
			if (opts.verbose) {
				console.log("SQL:", sql.includes("\n") ? "\n" + sql + " --> END SQL" : sql);
			}
			let prms = $`docker exec -it ${contId} \
/opt/mssql-tools/bin/sqlcmd -S localhost \
-U SA -P ${opts.password} -Q ${sql} || exit 1`;
			if (!opts.verbose && !opts.dryRun) {
				prms = prms?.quiet();
			}
			return prms;
		};

		console.log("Waiting for db to be up and running...");

		// https://docs.docker.com/engine/reference/run/#healthchecks
		await waitFor(() => sqlcmd("SELECT 1"));

		console.log("Creating the database and required data...");
		// TODO: parse the resposnse os SQL server and check for potential errors
		await sqlcmd(`CREATE DATABASE ${escapeId(opts.database)}`);

		vlog("Creating login");
		await sqlcmd(`CREATE LOGIN ${escapeUser(opts.user)} WITH PASSWORD = ${escapeStr(opts.password)}`);

		vlog("Creating user");
		await sqlcmd(`create user ${escapeUser(opts.user)} for login ${escapeUser(opts.user)}`, opts.database);

		// To check available roles: Select	[name] From sysusers Where issqlrole = 1
		vlog("Adding required permissions");
		await sqlcmd(`exec sp_addrolemember 'db_owner', ${escapeStr(opts.user)}`, opts.database);
	},
});
