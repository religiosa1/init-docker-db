export interface IDbCreateOptions {
	containerName: string;
	database: string;
	user: string;
	password: string;
	port: number;
	verbose?: boolean;
}

interface IDbCreator {
	name: string;
	port: number;
	defaultUser: string;
	create: (this: IDbCreator, opts: IDbCreateOptions) => Promise<void>;
}

export class DbCreator {
	readonly name: string;
	readonly port: number;
	readonly defaultUser: string;
	readonly create: (this: DbCreator, opts: IDbCreateOptions) => Promise<void>;

	constructor(opts: IDbCreator) {
		this.name = opts.name;
		this.port = opts.port;
		this.defaultUser = opts.defaultUser;
		this.create = function (args) {
			if (args.verbose) {
				console.log(this.name, "create", args);
			}
			return opts.create.call(this, args);
		};
	}
}
