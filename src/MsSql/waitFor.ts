import { setTimeout as pause } from "node:timers/promises";
type MaybePromise<T> = T | Promise<T>;

interface WaitForOpts {
	signal?: AbortSignal;
	preDelay?: number;
	minDelay?: number;
	maxDelay?: number;
	watchDogMs?: number;
	rate?: number;
}

/** truncated exponential backoff polling, until returns truthy */
export async function waitFor(
	exp: () => MaybePromise<unknown>,
	{ signal, watchDogMs = 60_000, minDelay = 200, maxDelay = 5_000, rate = 1.5, preDelay }: WaitForOpts = {}
): Promise<void> {
	if (!Number.isFinite(rate) || rate < 1.0 || rate > 10) {
		throw new RangeError(`Rate for exp backof must be in range 1 <= RATE <= 10, got ${rate}`);
	}
	if (!Number.isInteger(minDelay) || minDelay < 1) {
		throw new TypeError(`mindelay value must be a positive integer, got ${minDelay}`);
	}
	if (!Number.isInteger(maxDelay) || maxDelay < minDelay) {
		throw new TypeError(
			"maxDelay value must be a positive integer bigger than mindelay, " + `got ${maxDelay} (mindelay = ${minDelay})`
		);
	}
	if (!Number.isInteger(watchDogMs) || watchDogMs < maxDelay) {
		throw new TypeError(
			"watchDogMs value must be a positive integer bigger than maxDelay, " +
				`got ${watchDogMs} (maxDelay = ${maxDelay})`
		);
	}

	let currentDelay = minDelay;
	let nRuns = 0;
	if (preDelay) {
		await pause(preDelay, undefined, { signal });
	}
	await withTimeout(async () => {
		for (;;) {
			++nRuns;
			if (signal?.aborted) {
				throw signal.reason ?? new Error("waitFor was aborted without a reason");
			}
			// add ExecutionMax, and abort execution after it's passed, in case of long
			// hanging promises
			const result = await exp();
			if (result) break;
			await pause(currentDelay, undefined, { signal });
			currentDelay = Math.min(currentDelay * rate, maxDelay);
		}
	}, watchDogMs);
}

function withTimeout(cb: () => Promise<unknown>, timeout = 60_000): Promise<void> {
	return new Promise<void>((res, rej) => {
		const to = setTimeout(() => rej(new TimeoutError()), timeout);
		cb()
			.then(() => res())
			.catch(rej)
			.finally(() => {
				clearTimeout(to);
			});
	});
}

class TimeoutError extends Error {
	override name = "TimeoutError";
	constructor(msg?: string, opts?: ErrorOptions) {
		super(msg ?? "TimeoutError", opts);
	}
}
