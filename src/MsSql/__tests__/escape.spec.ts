import { describe, it, expect } from "bun:test";
import { escapeUser, escapeStr, escapeId } from "../escape";

describe("escape string", () => {
	it("always wraps a string in single quots", () => {
		expect(escapeStr("foo")).toBe("'foo'");
	});

	it("turns single quots in strin into a pair of singlequots", () => {
		expect(escapeStr("o'foo")).toBe("'o''foo'");
		expect(escapeStr("'ofoo")).toBe("'''ofoo'");
	});

	it("turns non-ascii chars into a CHAR() concatenation", () => {
		expect(escapeStr("foo\r\n")).toBe("'foo' + CHAR(13) + CHAR(10)");
	});

	it("still escapes concatenated pairs", () => {
		expect(escapeStr("foo\nb'ar")).toBe("'foo' + CHAR(10) + 'b''ar'");
	});

	it("doesn't append extra chars at start", () => {
		expect(escapeStr("\nfoo")).toBe("CHAR(10) + 'foo'");
	});

	it("encodes empty strings", () => {
		expect(escapeStr("")).toBe("''");
	});
});

describe("escape identifier", () => {
	it("always wrap identifier in square brackets", () => {
		expect(escapeId("foo")).toBe("[foo]");
		expect(escapeId("_foo")).toBe("[_foo]");
		expect(escapeId("foo.")).toBe("[foo.]");
	});

	it("treats characters from range [_$#@.\\-] as valid", () => {
		expect(escapeId("_$#.@-")).toBe("[_$#.@-]");
	});

	it("throw on non-printable characters in identifiers", () => {
		expect(() => escapeId("f\noo")).toThrow();
		expect(() => escapeId("f\boo")).toThrow();
	});

	it("throws on empty strings", () => {
		expect(() => escapeId("")).toThrow();
	});
});

describe("escape user", () => {
	it("throws on empty strings", () => {
		expect(() => escapeUser("")).toThrow();
		expect(() => escapeUser("a")).not.toThrow();
	});

	it("throws on non-printable characters", () => {
		expect(() => escapeUser("as\nda")).toThrow();
		expect(() => escapeUser("asda")).not.toThrow();
	});

	it("throws if user name gte 128 characters long", () => {
		expect(() => escapeUser("a".repeat(128))).toThrow();
		expect(() => escapeUser("a".repeat(127))).not.toThrow();
	});

	it("escapes user name in the same way as identifier", () => {
		expect(escapeUser("!asd")).toBe(escapeId("!asd"));
		expect(escapeUser("_$#.@-")).toBe(escapeId("_$#.@-"));
		expect(escapeUser("as.d")).toBe(escapeId("as.d"));
	});
});
