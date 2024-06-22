import yargs from "yargs/yargs";
import { hideBin } from "yargs/helpers";

import generateName from "boring-name-generator";

import { LazyReadline, ReadlineDisabledError } from "./src/LazyReadline";
import { PossibleDbTypes, isValidDbType } from "./src/PossibleDbTypes";
import { DbCreator, type IDbCreateOptions } from "./src/DbCreator";
import { Postgres } from "./src/Postgres";
import { MsSql } from "./src/MsSql";
import { MySql } from "./src/MySql";

const creators: Readonly<Record<PossibleDbTypes, DbCreator>> = {
	postgres: Postgres,
	mssql: MsSql,
	mysql: MySql,
};

const args = await yargs(hideBin(process.argv))
	.usage("Create a database container")
	.positional("name", {
		type: "string",
		describe: "name of the database container to be created",
	})
	.option("type", {
		alias: "t",
		type: "string",
		choices: PossibleDbTypes,
		description: "database type",
	})
	.option("user", {
		alias: "u",
		type: "string",
		description: "database user",
	})
	.option("database", {
		alias: "d",
		type: "string",
		description: "database name",
	})
	.option("password", {
		alias: "p",
		type: "string",
		description: "user's password",
	})
	.option("port", {
		alias: "P",
		type: "number",
		description: "TCP port to which database will be mapped to",
	})
	.option("tag", {
		alias: "T",
		type: "string",
		description: "docker tag to use with the container",
	})
	.option("non-interactive", {
		alias: "n",
		type: "boolean",
		description: "exit, if some of the required params are missing",
	})
	.option("verbose", {
		alias: "v",
		type: "boolean",
		description: "Run with verbose logging",
	})
	.help("h")
	.alias("h", "help")
	.parse();
type CliArgs = typeof args;

await main(args);

async function main(args: CliArgs): Promise<void> {
	using rl = new LazyReadline(!!args.nonInteractive);
	try {
		var creator = await getCreator(rl, args.type);
		var options = await getOptions(rl, creator, args);
	} catch (e) {
		if (e instanceof ReadlineDisabledError) {
			console.log("Some of required data is missing or incorrect and non-interactive flag is provided, exiting");
			process.exit(1);
		}
		throw e;
	}
	await creator.create(options);
	console.log("Done");
}

async function getCreator(rl: LazyReadline, type: string | undefined): Promise<DbCreator> {
	let creator: DbCreator | undefined = creators[type as PossibleDbTypes];
	while (!creator) {
		const TYPE_DEFAULT = "postgres";
		const answer = await rl.question(`database type? [${PossibleDbTypes.join()}] (${TYPE_DEFAULT}): `);
		if (!answer.trim()) {
			creator = creators[TYPE_DEFAULT];
			break;
		}
		if (isValidDbType(answer)) {
			creator = creators[answer];
		} else {
			console.log("Incorrect DB type, possible types are: " + PossibleDbTypes.join());
		}
	}
	return creator;
}

async function getOptions(rl: LazyReadline, creator: DbCreator, args: CliArgs): Promise<IDbCreateOptions> {
	const opts: IDbCreateOptions = {
		database: args.database!,
		user: args.user!,
		password: args.password!,
		containerName: args.name!,
		port: args.port || creator.port,
		tag: args.tag || creator.defaultTag,
		verbose: args.verbose,
	};

	if (!opts.database) {
		const DATABASE_DEFAULT = "db";
		opts.database = await rl.question(`database name? (${DATABASE_DEFAULT}): `);
		opts.database ||= DATABASE_DEFAULT;
	}
	if (!opts.user) {
		const USER_DEFAULT = creator.defaultUser;
		opts.user = await rl.question(`database user? (${USER_DEFAULT}): `);
		opts.user ||= USER_DEFAULT;
	}
	if (!opts.password) {
		for (;;) {
			const PASSWORD_DEFAULT = creator.defaultPassword;
			opts.password = await rl.question(`database password? (${PASSWORD_DEFAULT}): `);
			opts.password ||= PASSWORD_DEFAULT;

			const [valid, message] = creator.isPasswordValid(opts.password);
			if (valid) {
				break;
			}
			console.log(message || "Password is invalid");
		}
	}
	if (!opts.containerName) {
		const name = generateName().dashed;
		opts.containerName = await rl.question(`docker container name? (${name}): `);
		opts.containerName ||= name;
	}
	return opts;
}
