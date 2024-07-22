export const changeFirstLetterToUpperCase = (letter: string) => {
  if (letter?.split(" ")[1] !== undefined) {
    const a0 = letter?.split(" ")[0]?.substring(0, 1)?.toUpperCase();
    const b0 = letter?.split(" ")[0]?.substring(1);

    const a1 = letter?.split(" ")[1]?.substring(0, 1)?.toUpperCase();
    const b1 = letter?.split(" ")[1]?.substring(1);

    return a0 + b0 + " " + (a1 + b1);
  }
  const a = letter?.substring(0, 1)?.toUpperCase();
  const b = letter?.substring(1);
  return a + b;
};
