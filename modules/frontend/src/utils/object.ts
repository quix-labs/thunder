export function flatObject(obj: { [key: string]: any }, target?: { [key: string]: any }, prefix?: string) {
    target = target || {};
    prefix = prefix || "";

    Object.keys(obj).forEach(function (key) {
        if (typeof (obj[key]) === "object" && obj[key] !== null) {
            flatObject(obj[key], target, prefix + key + ".");
        } else {
            return target[prefix + key] = obj[key];
        }
    });

    return target;
}