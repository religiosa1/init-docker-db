import { DbCreator } from "./DbCreator";

export const MySql = new DbCreator({
	name: "mysql",
	port: 3306,
	defaultUser: "mysql",
	async create(opts) {
		throw new Error("Not implemented");
	},
});
