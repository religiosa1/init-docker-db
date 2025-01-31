import { DbCreator } from "./DbCreator";

export const Mongo = new DbCreator({
	name: "mongo",
	port: 27017,
	defaultUser: "mongo",
	async create($, opts) {
		// https://hub.docker.com/_/mongo
		await $`docker run --name ${opts.containerName} \
-e MONGO_INITDB_ROOT_PASSWORD=${opts.password} \
-e MONGO_INITDB_ROOT_USERNAME=${opts.user} \
-e MONGO_INITDB_DATABASE=${opts.database} \
-p ${this.port}:${opts.port} \
-d mongo:${opts.tag}`;
	},
});
