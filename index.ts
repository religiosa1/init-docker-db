import yargs from "yargs/yargs";
import { hideBin } from "yargs/helpers";
import { version } from "./package.json" with { type: "json" };

import generateName from "boring-name-generator";

import { LazyReadline, ReadlineDisabledError } from "./src/LazyReadline";
import { dbTypesList, type DbType, isValidDbType, dbTypes } from "./src/dbTypes";
import { DbCreator, type IDbCreateOptions } from "./src/DbCreator";

class ValidationError extends Error {}

const args = await yargs(hideBin(process.argv))
	.usage("Create a disposable database docker container")
	.positional("containerName", {
		type: "string",
		array: false,
		describe: "name of the database container to be created",
	})
	.option("type", {
		alias: "t",
		type: "string",
		choices: dbTypesList,
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
	.version(version)
	.alias("h", "help")
	.example([
		// Not using $0 yargs interpolation, as it will resole the name to "bun"
		["init-docker-db", "Run in wizard mode"],
		["init-docker-db -t mssql -u app_user", "Create a MsSQL database using provided username"],
	])
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
		if (e instanceof ValidationError) {
			console.log("Options validation error:", e.message);
			process.exit(2);
		}
		throw e;
	}
	await creator.create(options);
	console.log("Done");
}

async function getCreator(rl: LazyReadline, type: string | undefined): Promise<DbCreator> {
	let creator: DbCreator | undefined = dbTypes[type as DbType];
	while (!creator) {
		const TYPE_DEFAULT = "postgres";
		const answer = await rl.question(`database type? [${dbTypesList.join()}] (${TYPE_DEFAULT}): `);
		if (!answer.trim()) {
			creator = dbTypes[TYPE_DEFAULT];
			break;
		}
		if (isValidDbType(answer)) {
			creator = dbTypes[answer];
		} else {
			console.log("Incorrect DB type, possible types are: " + dbTypesList.join());
		}
	}
	return creator;
}

async function getOptions(rl: LazyReadline, creator: DbCreator, args: CliArgs): Promise<IDbCreateOptions> {
	const opts: IDbCreateOptions = {
		database: args.database!,
		user: args.user!,
		password: args.password!,
		containerName: args["_"]?.[0]?.toString(),
		port: args.port || creator.port,
		tag: args.tag || creator.defaultTag,
		verbose: args.verbose,
	};

	// validating existing password first if it's there
	if (opts.password) {
		const [valid, message] = creator.isPasswordValid(opts.password);
		if (!valid) {
			throw new ValidationError(message || "Provided password does not meet the requirements");
		}
	}

	// Filling out the rest of missing data in interactive mode
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
			console.log(message || "Provided password does not meet the requirements");
		}
	}
	if (!opts.containerName) {
		const name = generateName().dashed;
		opts.containerName = await rl.question(`docker container name? (${name}): `);
		opts.containerName ||= name;
	}
	return opts;
}
