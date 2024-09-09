export function throttle<T extends (...args: any[]) => Promise<any>>(func: T, limit: number): (...args: Parameters<T>) => Promise<void> {
    let inThrottle: boolean;
    let lastArgs: Parameters<T> | null = null;

    return async function(...args: Parameters<T>) {
        if (inThrottle) {
            lastArgs = args;
        } else {
            inThrottle = true;
            await func(...args);
            setTimeout(async () => {
                inThrottle = false;
                if (lastArgs) {
                    await func(...lastArgs);
                    lastArgs = null;
                }
            }, limit);
        }
    };
}