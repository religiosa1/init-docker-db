import { DbCreator, type PasswordValidityTuple } from "../DbCreator";
import { createVerboseShell } from "../createVerboseShell";

import { MsSqlPwdValidityEnum, msSqlPwdValidityEnumMessage } from "./MsSqlPwdValidityEnum";
import { escapeId, escapeUser, escapeStr } from "./escape";

const PWD_COMPLEXITY_REGEX = /^(?=.*[A-Z])(?=.*[a-z])(?=.*\d).+$/;

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
		if (!PWD_COMPLEXITY_REGEX.test(password)) {
			return [false, msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.TooSimple]];
		}
		return [true, msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.Valid]];
	},

	async create(opts) {
		const vlog = (...args: unknown[]) => (opts.verbose ? console.log("verbose:", ...args) : () => {});

		const $ = createVerboseShell(opts.verbose);

		// https://mcr.microsoft.com/product/mssql/server/about
		const shellOutput = await $`docker run -e ACCEPT_EULA=Y \
--name ${opts.containerName}\
--hostname ${opts.containerName}\
-e MSSQL_SA_PASSWORD=${opts.password}\
-p ${this.port}:${opts.port}\
-d mcr.microsoft.com/mssql/server:${this.defaultTag}`;
		const contId = shellOutput.text().trim();

		const sqlcmd = async (sql: string) => $`docker exec -it ${contId} \
/opt/mssql-tools/bin/sqlcmd -S localhost \
-U SA -P Password12 -Q ${sql}`;

		vlog("Waiting for db to be up and running");

		// TODO: Replace with healthcheck polling
		// https://docs.docker.com/engine/reference/run/#healthchecks
		await pause(20_000);

		vlog("Creating the database");
		// TODO: parse the resposnse os SQL server and check for potential errors
		await sqlcmd(`CREATE DATABASE ${escapeId(opts.database)}`);

		vlog("Creating login");
		await sqlcmd(`CREATE LOGIN ${escapeUser(opts.user)} WITH PASSWORD = ${escapeStr(opts.password)}`);

		vlog("Creating user");
		await sqlcmd(
			`use ${escapeId(opts.database)}\n` + //
				`create user ${escapeUser(opts.user)} for login ${escapeUser(opts.user)}`
		);

		vlog("Adding required permissions");
		await sqlcmd(`exec sp_addrolemember ${escapeStr(opts.user)}, 'dbowner'`);
	},
});

function pause(timeout = 3000) {
	return new Promise<void>((res) => setTimeout(res, timeout));
}
