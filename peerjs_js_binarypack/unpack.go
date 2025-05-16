package peerjs_js_binarypack

import (
	"github.com/dop251/goja"
)

func UnpackStr(bin []byte) string {
	vm := goja.New()

	// JavaScript 코드 실행
	_, err := vm.RunString(`
        var c = class {
                constructor(arr) {
                    this.index = 0, this.dataBuffer = arr.buffer, this.dataView = arr, this.length = this.dataBuffer.byteLength
                }
                unpack() {
                    var e = this.unpack_uint8();
                    if (e < 128) return e;
                    if ((e ^ 224) < 32) return (e ^ 224) - 32;
                    let t;
                    if ((t = e ^ 160) <= 15) return this.unpack_raw(t);
                    if ((t = e ^ 176) <= 15) return this.unpack_string(t);
                    if ((t = e ^ 144) <= 15) return this.unpack_array(t);
                    if ((t = e ^ 128) <= 15) return this.unpack_map(t);
                    switch (e) {
                        case 192:
                            return null;
                        case 193:
                            return;
                        case 194:
                            return !1;
                        case 195:
                            return !0;
                        case 202:
                            return this.unpack_float();
                        case 203:
                            return this.unpack_double();
                        case 204:
                            return this.unpack_uint8();
                        case 205:
                            return this.unpack_uint16();
                        case 206:
                            return this.unpack_uint32();
                        case 207:
                            return this.unpack_uint64();
                        case 208:
                            return this.unpack_int8();
                        case 209:
                            return this.unpack_int16();
                        case 210:
                            return this.unpack_int32();
                        case 211:
                            return this.unpack_int64();
                        case 212:
                            return;
                        case 213:
                            return;
                        case 214:
                            return;
                        case 215:
                            return;
                        case 216:
                            t = this.unpack_uint16(); return this.unpack_string(t);
                        case 217:
                            t = this.unpack_uint32(); return this.unpack_string(t);
                        case 218:
                            t = this.unpack_uint16(); return this.unpack_raw(t);
                        case 219:
                            t = this.unpack_uint32(); return this.unpack_raw(t);
                        case 220:
                            t = this.unpack_uint16(); return this.unpack_array(t);
                        case 221:
                            t = this.unpack_uint32(); return this.unpack_array(t);
                        case 222:
                            t = this.unpack_uint16(); return this.unpack_map(t);
                        case 223:
                            t = this.unpack_uint32(); return this.unpack_map(t);
                    }
                }
                unpack_uint8() {
                    var e = this.dataView[this.index] & 255;
                    this.index++; return e;
                }
                unpack_uint16() {
                    var e = this.read(2);
                    var t = (e[0] & 255) * 256 + (e[1] & 255);
                    this.index += 2; return t;
                }
                unpack_uint32() {
                    var e = this.read(4);
                    var t = ((e[0] * 256 + e[1]) * 256 + e[2]) * 256 + e[3];
                    this.index += 4; return t;
                }
                unpack_uint64() {
                    var e = this.read(8);
                    var t = ((((((e[0] * 256 + e[1]) * 256 + e[2]) * 256 + e[3]) * 256 + e[4]) * 256 + e[5]) * 256 + e[6]) * 256 + e[7];
                    this.index += 8; return t;
                }
                unpack_int8() {
                    var e = this.unpack_uint8();
                    return e < 128 ? e : e - 256;
                }
                unpack_int16() {
                    var e = this.unpack_uint16();
                    return e < 32768 ? e : e - 65536;
                }
                unpack_int32() {
                    var e = this.unpack_uint32();
                    return e < 2 ** 31 ? e : e - 2 ** 32;
                }
                unpack_int64() {
                    var e = this.unpack_uint64();
                    return e < 2 ** 63 ? e : e - 2 ** 64;
                }
                unpack_raw(e) {
                    if (this.length < this.index + e) throw new Error("BinaryPackFailure");
                    var t = this.dataBuffer.slice(this.index, this.index + e);
                    this.index += e; return t;
                }
                unpack_string(e) {
                    var t = this.read(e);
                    var f = 0;
                    var i = "";
                    let r, n;
                    for (; f < e;) r = t[f], r < 160 ? (n = r, f++) : (r ^ 192) < 32 ? (n = (r & 31) << 6 | t[f + 1] & 63, f += 2) : (r ^ 224) < 16 ? (n = (r & 15) << 12 | (t[f + 1] & 63) << 6 | t[f + 2] & 63, f += 3) : (n = (r & 7) << 18 | (t[f + 1] & 63) << 12 | (t[f + 2] & 63) << 6 | t[f + 3] & 63, f += 4), i += String.fromCodePoint(n);
                    this.index += e; return i;
                }
                unpack_array(e) {
                    var t = new Array(e);
                    for (let f = 0; f < e; f++) t[f] = this.unpack();
                    return t;
                }
                unpack_map(e) {
                    let t = {};
                    for (let f = 0; f < e; f++) {
                        let i = this.unpack();
                        t[i] = this.unpack();
                    }
                    return t;
                }
                unpack_float() {
                    var e = this.unpack_uint32();
                    var t = e >> 31;
                    var f = (e >> 23 & 255) - 127;
                    var i = e & 8388607 | 8388608;
                    return (t === 0 ? 1 : -1) * i * 2 ** (f - 23);
                }
                unpack_double() {
                    var e = this.unpack_uint32();
                    var t = this.unpack_uint32();
                    var f = e >> 31;
                    var i = (e >> 20 & 2047) - 1023;
                    var n = (e & 1048575 | 1048576) * 2 ** (i - 20) + t * 2 ** (i - 52);
                    return (f === 0 ? 1 : -1) * n;
                }
                read(e) {
                    var t = this.index;
                    if (t + e <= this.length) return this.dataView.subarray(t, t + e);
                    throw new Error("BinaryPackFailure: read index out of range");
                }
            };
    `)
	if err != nil {
		panic(err)
	}
	_, err = vm.RunString("var run = (pp) => new c(new Uint8Array(pp)).unpack();")
	if err != nil {
		panic(err)
	}
	processFunc, ok := goja.AssertFunction(vm.Get("run"))
	if !ok {
		panic("processData is not a function")
	}
	result, err := processFunc(goja.Undefined(), vm.ToValue(bin))
	if err != nil {
		panic(err)
	}
	return result.Export().(string)
}
