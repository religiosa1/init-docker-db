import { $ } from "bun";
import { DbCreator } from "./DbCreator";

export const Postgres = new DbCreator({
	name: "postgres",
	port: 5432,
	defaultUser: "postgres",
	async create(opts) {
		await $`docker run --name ${opts.containerName} \
-e POSTGRES_PASSWORD=${opts.password} \
-e POSTGRES_USER=${opts.user} \
-e POSTGRES_DB=${opts.database} \
-p ${this.port}:${opts.port} \
-d postgres`;
	},
});
