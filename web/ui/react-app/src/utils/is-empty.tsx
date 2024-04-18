/**
 * Returns whether the array is empty, null, or undefined
 *
 * @param arg - The array to check
 * @returns Whether the array is empty, null, or undefined
 */
export const isEmptyArray = <T extends unknown[] | undefined>(
  arg: T
): boolean => ((arg as unknown[]) ?? []).length === 0;

export const isEmptyObject = <T extends Record<string, unknown> | undefined>(
  arg: T
): boolean => Object.keys(arg ?? {}).length === 0;
