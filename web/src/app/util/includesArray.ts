export function includesArray(arrA: string[], arrB: string[]): boolean {
  if (arrA.length < arrB.length) {
    return false;
  }

  for (let i = 0; i < arrB.length; i++) {
    if (arrA[i] !== arrB[i]) {
      return false;
    }
  }

  return true;
}
