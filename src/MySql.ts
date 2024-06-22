import { DbCreator } from "./DbCreator";
import { createVerboseShell } from "./createVerboseShell";

export const MySql = new DbCreator({
	name: "mysql",
	port: 3306,
	defaultUser: "mysql",
	defaultTag: "lts",
	async create(opts) {
		const $ = createVerboseShell(opts.verbose);
		// https://hub.docker.com/_/mysql
		await $`docker run --name ${opts.containerName}\
-e MYSQL_USER=${opts.password}\
-e MYSQL_ROOT_PASSWORD=${opts.password}\
-e MYSQL_PASSWORD=${opts.password}\
-e MYSQL_DATABASE=${opts.database}\
-p ${this.port}:${opts.port}\
-d mysql:${this.defaultTag}`;
	},
});
