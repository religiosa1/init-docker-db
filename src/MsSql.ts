import { DbCreator } from "./DbCreator";
import { createVerboseShell } from "./createVerboseShell";

export const MsSql = new DbCreator({
	name: "mssql",
	port: 1433,
	defaultUser: "mssql",
	defaultTag: "2022-latest",
	// TODO password validation, as group policy requires for it to be 10 chars long and have at
	// least one upper and lowercase and a digit
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
		// TODO: escaping of sql commands
		// TODO: parse the resposnse os SQL server and check for potential errors
		await sqlcmd(`CREATE DATABASE ${opts.database}`);

		vlog("Creating login");
		await sqlcmd(`CREATE LOGIN ${opts.user} WITH PASSWORD = '${opts.password}'`);

		vlog("Creating user");
		await sqlcmd(
			`use ${opts.database}
create user ${opts.user} for login ${opts.user}
`
		);

		vlog("Adding required permissions");
		await sqlcmd(`exec sp_addrolemember '${opts.user}', 'dbowner'`);
	},
});

function pause(timeout = 3000) {
	return new Promise<void>((res) => setTimeout(res, timeout));
}
