type MaybePromise<T> = T | Promise<T>;

interface WaitForOpts {
	signal?: AbortSignal;
	watchDogMs?: number;
	initialDelay?: number;
	maxDelay?: number;
	rate?: number;
}

/** truncated exponential backoff polling, until doesn't throw */
export async function waitFor(
	exp: () => MaybePromise<unknown>,
	{ signal, watchDogMs = 60_000, initialDelay = 200, maxDelay = 5_000, rate = 1.5 }: WaitForOpts = {}
): Promise<void> {
	if (!Number.isFinite(rate) || rate < 1.0 || rate > 10) {
		throw new RangeError(`Rate for exp backof cannot mesut be in range 1 <= RATE <= 10, got ${rate}`);
	}
	if (!Number.isInteger(initialDelay) || initialDelay < 1) {
		throw new TypeError(`Initial delay value must be a positive integer, got ${initialDelay}`);
	}
	if (!Number.isInteger(maxDelay) || maxDelay < initialDelay) {
		throw new TypeError(
			"Max delay value must be a positive integer bigger than initialDelay, " +
				`got ${maxDelay} (initialDelat = ${initialDelay})`
		);
	}
	if (!Number.isInteger(watchDogMs) || watchDogMs < maxDelay) {
		throw new TypeError(
			"watchDogMs value must be a positive integer bigger than maxDelay, " +
				`got ${watchDogMs} (maxDelay = ${maxDelay})`
		);
	}

	let currentDelay = initialDelay;
	let nRuns = 0;
	await withTimeout(async () => {
		for (;;) {
			++nRuns;
			if (signal?.aborted) {
				throw signal.reason ?? new Error("waitFor aborted without a reason");
			}
			try {
				await exp();
				break;
			} catch {}
			await pause(currentDelay, { signal });
			currentDelay = Math.max(currentDelay * rate, maxDelay);
		}
	}, watchDogMs);
}

function pause(timeout: number, { signal }: { signal?: AbortSignal } = {}): Promise<void> {
	return new Promise<void>((res, rej) => {
		const to = setTimeout(() => {
			signal?.removeEventListener("abort", handleAbort);
			res();
		}, timeout);
		signal?.addEventListener("abort", handleAbort, { once: true });
		function handleAbort() {
			clearTimeout(to);
			rej(signal?.reason ?? new Error("aborted without a reason"));
		}
	});
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
