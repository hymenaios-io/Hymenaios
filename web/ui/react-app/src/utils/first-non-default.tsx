import isEmptyOrNull from './is-empty-or-null';

/**
 * The first non-empty value.
 *
 * - All undefined/empty = ""
 *
 * @param args - The list of arguments to check.
 * @returns The first non-empty value from the list of arguments, as a string.
 */
const firstNonDefault: (...args: unknown[]) => string = (
	...args: unknown[]
) => {
	// Find the first non-empty argument.
	for (const arg of args) {
		if (!isEmptyOrNull(arg)) return `${arg}`;
	}
	// No non-empty argument found, return an empty string.
	return '';
};

export default firstNonDefault;
