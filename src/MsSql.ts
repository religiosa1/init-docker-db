import { DbCreator } from "./DbCreator";

export const MsSql = new DbCreator({
	name: "mssql",
	port: 1433,
	defaultUser: "mssql",
	async create(opts) {
		throw new Error("Not implemented");
	},
});
