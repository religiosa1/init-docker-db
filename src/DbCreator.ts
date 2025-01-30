export interface IDbCreateOptions {
	dryRun?: boolean;
	containerName: string;
	database: string;
	user: string;
	password: string;
	port: number;
	tag: string;
	verbose?: boolean;
}

export type PasswordValidityTuple = [status: true, message?: string] | [status: false, message: string];

interface IDbCreator {
	name: string;
	port: number;
	defaultUser: string;
	defaultTag: string;
	defaultPassword: string;
	create: (this: IDbCreator, opts: IDbCreateOptions) => Promise<void>;
	isPasswordValid(password: string): PasswordValidityTuple;
}

type WithDefaults<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>;

export class DbCreator implements IDbCreator {
	readonly name: string;
	readonly port: number;
	readonly defaultUser: string;
	readonly defaultTag: string;
	readonly defaultPassword: string;
	readonly create: (this: IDbCreator, opts: IDbCreateOptions) => Promise<void>;

	constructor(opts: WithDefaults<IDbCreator, "defaultPassword" | "isPasswordValid" | "defaultTag">) {
		this.name = opts.name;
		this.port = opts.port;
		this.defaultUser = opts.defaultUser;
		this.defaultTag = opts.defaultTag || "latest";
		this.defaultPassword = opts.defaultPassword || "123456";

		if (opts.isPasswordValid) {
			this.isPasswordValid = opts.isPasswordValid;
		}

		this.create = function (args) {
			if (args.verbose) {
				console.log(this.name, "create", args);
			}
			return opts.create.call(this, args);
		};
	}

	isPasswordValid(_password: string): PasswordValidityTuple {
		return [true];
	}
}
