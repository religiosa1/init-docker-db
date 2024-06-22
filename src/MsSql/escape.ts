export function escapeUser(name: string): string {
	if (name.length >= 128) {
		throw new Error("User name can be up to 128 charaters long");
	}
	return escapeId(name);
}

export function escapeId(name: string): string {
	if (!name) {
		throw new Error("MsSQL identifier cannot be empty");
	}
	if (name.includes("[") || name.includes("]")) {
		throw new Error("MsSQL identifier cannot contain '[' or ']' character");
	}
	if (!isPrintable(name)) {
		throw new Error("Identifiers cannot contain non-printable characters");
	}
	return `[${name}]`;
}

export function escapeStr(str: string): string {
	if (!str) {
		return "''";
	}

	const spans = splitString(str);
	const segments = spans.map(([value, type]) =>
		type === "printable" ? escapePrintableStr(value) : escapeNonPrintableStr(value)
	);

	return segments.join(" + ");
}

function escapePrintableStr(str: string): string {
	return "'" + str.replaceAll("'", "''") + "'";
}

function escapeNonPrintableStr(char: string): string {
	return `CHAR(${char.charCodeAt(0)})`;
}

type Span = [value: string, type: "printable" | "nonprintable"];
function splitString(str: string): Span[] {
	let result: Span[] = [];
	let lastNonPrintableIndex = -1;
	for (let i = 0; i < str.length; i++) {
		const charCode = str.charCodeAt(i);
		if (isPrintable(charCode)) {
			continue;
		}
		if (i !== 0 && lastNonPrintableIndex !== i - 1) {
			result.push([str.substring(lastNonPrintableIndex + 1, i), "printable"]);
		}
		result.push([str.substring(i, i + 1), "nonprintable"]);
		lastNonPrintableIndex = i;
	}
	if (lastNonPrintableIndex !== str.length - 1) {
		result.push([str.substring(lastNonPrintableIndex + 1), "printable"]);
	}
	return result;
}

function isPrintable(str: string): boolean;
function isPrintable(charCode: number): boolean;
function isPrintable(charCodeOrString: number | string): boolean {
	if (typeof charCodeOrString === "number") {
		const charCode = charCodeOrString;
		return charCode >= 32 && charCode < 127;
	}
	const str = charCodeOrString;
	for (let i = 0; i < str.length; i++) {
		const charCode = str.charCodeAt(i);
		if (charCode < 32 || charCode >= 127) {
			return false;
		}
	}
	return true;
}
