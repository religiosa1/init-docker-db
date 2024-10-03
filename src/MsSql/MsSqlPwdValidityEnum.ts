export enum MsSqlPwdValidityEnum {
	Valid,
	Empty,
	TooShort,
	TooSimple,
}

export const msSqlPwdValidityEnumMessage: Record<MsSqlPwdValidityEnum, string> = {
	[MsSqlPwdValidityEnum.Valid]: "Valid",
	[MsSqlPwdValidityEnum.Empty]: "Password can't be empty",
	[MsSqlPwdValidityEnum.TooShort]: "Password is too short (must be at least 10 chars)",
	[MsSqlPwdValidityEnum.TooSimple]:
		"Password doesn't meet the complexity requirements " +
		"(must contain 3 out of 4 char types: lowercase char, uppercase char, digit, nonalphanumeric)",
};
