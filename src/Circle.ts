export default class Circle {
    constructor(readonly radius: number) {}

    getArea(): number {
        return Math.PI * this.radius * this.radius;
    }
}