import { test, expect, describe } from "bun:test";
import { MsSql } from "../MsSql";
import { MsSqlPwdValidityEnum, msSqlPwdValidityEnumMessage } from "../MsSqlPwdValidityEnum";

describe("Password validation", () => {
	const valid = msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.Valid];

	test("default password is ok", () => {
		expect(MsSql.isPasswordValid(MsSql.defaultPassword)).toEqual([true, valid]);
	});

	test("empty password", () => {
		expect(MsSql.isPasswordValid("")).toEqual([false, msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.Empty]]);
	});

	test("password too short", () => {
		expect(MsSql.isPasswordValid("Password12")).toEqual([true, valid]);
		expect(MsSql.isPasswordValid("Password1")).toEqual([
			false,
			msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.TooShort],
		]);
	});

	test("password too simple", () => {
		expect(MsSql.isPasswordValid("PASSWORD1!")).toEqual([true, valid]);
		expect(MsSql.isPasswordValid("password1!")).toEqual([true, valid]);
		expect(MsSql.isPasswordValid("password12")).toEqual([
			false,
			msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.TooSimple],
		]);
		expect(MsSql.isPasswordValid("PASSWORD12")).toEqual([
			false,
			msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.TooSimple],
		]);
		expect(MsSql.isPasswordValid("PASSWORD!!")).toEqual([
			false,
			msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.TooSimple],
		]);
		expect(MsSql.isPasswordValid("0123456789")).toEqual([
			false,
			msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.TooSimple],
		]);
	});

	test("complexity doesn't depend on order", () => {
		expect(MsSql.isPasswordValid("pAssword12")).toEqual([
			true,
			msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.Valid],
		]);
		expect(MsSql.isPasswordValid("12Password")).toEqual([
			true,
			msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.Valid],
		]);
		expect(MsSql.isPasswordValid("Pass12word")).toEqual([
			true,
			msSqlPwdValidityEnumMessage[MsSqlPwdValidityEnum.Valid],
		]);
	});
});
