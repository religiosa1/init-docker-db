import { $ } from "bun";
import { DbCreator } from "./DbCreator";

export const MySql = new DbCreator({
	name: "mysql",
	port: 3306,
	defaultUser: "mysql",
	async create(opts) {
		await $`docker run --name ${opts.containerName}\
-e MYSQL_USER=${opts.password}\
-e MYSQL_ROOT_PASSWORD=${opts.password}\
-e MYSQL_PASSWORD=${opts.password}\
-e MYSQL_DATABASE=${opts.database}\
-p ${this.port}:${opts.port}\
-d mysql:lts`;
	},
});
