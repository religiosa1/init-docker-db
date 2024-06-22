export const PossibleDbTypes = ["postgres", "mysql"] as const;
export type PossibleDbTypes = (typeof PossibleDbTypes)[number];
export function isValidDbType(value: any): value is PossibleDbTypes {
	return PossibleDbTypes.includes(value);
}
