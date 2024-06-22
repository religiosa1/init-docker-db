import { Postgres } from "./Postgres";
import { MsSql } from "./MsSql";
import { MySql } from "./MySql";
import { Mongo } from "./Mongo";

export const dbTypes = Object.freeze({
	postgres: Postgres,
	mysql: MySql,
	mssql: MsSql,
	mongo: Mongo,
});

export const dbTypesList = Array.from(Object.keys(dbTypes)) as DbType[];
export type DbType = keyof typeof dbTypes;

export function isValidDbType(k: string): k is DbType {
	return Object.hasOwn(dbTypes, k);
}
