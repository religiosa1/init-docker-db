import { DbCreator } from "./DbCreator";

export const MsSql = new DbCreator({
	name: "mssql",
	port: 1433,
	defaultUser: "mssql",
	defaultTag: "2022-latest",
	async create(opts) {
		throw new Error("Not implemented");
	},
});
