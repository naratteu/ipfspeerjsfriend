package peerjs_js_binarypack

import (
	"github.com/dop251/goja"
)

func PackStr(str string) []byte {
	vm := goja.New()

	_, err := vm.RunString(`
        class TextEncoder {
            encode(str) {
                let utf8 = [];
                for (let i = 0; i < str.length; i++) {
                    let code = str.charCodeAt(i);
                    if (code < 128) {
                        utf8.push(code);
                    } else if (code < 2048) {
                        utf8.push(192 | (code >> 6));
                        utf8.push(128 | (code & 63));
                    } else {
                        utf8.push(224 | (code >> 12));
                        utf8.push(128 | ((code >> 6) & 63));
                        utf8.push(128 | (code & 63));
                    }
                }
                return new Uint8Array(utf8);
            }
        }
    `)
	if err != nil {
		panic(err)
	}

	// JavaScript 코드 실행
	_, err = vm.RunString(`
            var a = class {
                constructor() {
                    this._parts = [];
                }
                append_buffer(e) {
                    this._parts.push(...e);
                }
                append(e) {
                    this._parts.push(e);
                }
            };
            let p = class {
                pack(e) {
                    if (typeof e == "string") this.pack_string(e);
                    else if (typeof e == "number") Math.floor(e) === e ? this.pack_integer(e) : this.pack_double(e);
                    else if (typeof e == "boolean") e === !0 ? this._bufferBuilder.append(195) : e === !1 && this._bufferBuilder.append(194);
                    else if (e === void 0) this._bufferBuilder.append(192);
                    else if (typeof e == "object")
                        if (e === null) this._bufferBuilder.append(192);
                        else {
                            let t = e.constructor;
                            if (e instanceof Array) {
                                let f = this.pack_array(e);
                                if (f instanceof Promise) return f.then(() => {})
                            } else if (e instanceof ArrayBuffer) this.pack_bin(new Uint8Array(e));
                            else if ("BYTES_PER_ELEMENT" in e) {
                                let f = e;
                                this.pack_bin(new Uint8Array(f.buffer, f.byteOffset, f.byteLength))
                            } else if (e instanceof Date) this.pack_string(e.toString());
                            else {
                                // if (e instanceof Blob) return e.arrayBuffer().then(f => {
                                //     this.pack_bin(new Uint8Array(f));
                                // });
                                if (t == Object || t.toString().startsWith("class")) {
                                    let f = this.pack_object(e);
                                    if (f instanceof Promise) return f.then(() => {})
                                } else throw new Error("Type not yet supported")
                            }
                        }
                    else throw new Error("Type not yet supported");
                }
                pack_bin(e) {
                    var t = e.length;
                    if (t <= 15) this.pack_uint8(160 + t);
                    else if (t <= 65535) this._bufferBuilder.append(218), this.pack_uint16(t);
                    else if (t <= 4294967295) this._bufferBuilder.append(219), this.pack_uint32(t);
                    else throw new Error("Invalid length");
                    this._bufferBuilder.append_buffer(e);
                }
                pack_string(e) {
                    var t = new TextEncoder().encode(e);
                    var f = t.length;
                    if (f <= 15) this.pack_uint8(176 + f);
                    else if (f <= 65535) this._bufferBuilder.append(216), this.pack_uint16(f);
                    else if (f <= 4294967295) this._bufferBuilder.append(217), this.pack_uint32(f);
                    else throw new Error("Invalid length");
                    this._bufferBuilder.append_buffer(t);
                }
                pack_array(e) {
                    var t = e.length;
                    if (t <= 15) this.pack_uint8(144 + t);
                    else if (t <= 65535) this._bufferBuilder.append(220), this.pack_uint16(t);
                    else if (t <= 4294967295) this._bufferBuilder.append(221), this.pack_uint32(t);
                    else throw new Error("Invalid length");
                    var f = i => {
                        if (i < t) {
                            let r = this.pack(e[i]);
                            return r instanceof Promise ? r.then(() => f(i + 1)) : f(i + 1);
                        }
                    };
                    return f(0);
                }
                pack_integer(e) {
                    if (e >= -32 && e <= 127) this._bufferBuilder.append(e & 255);
                    else if (e >= 0 && e <= 255) this._bufferBuilder.append(204), this.pack_uint8(e);
                    else if (e >= -128 && e <= 127) this._bufferBuilder.append(208), this.pack_int8(e);
                    else if (e >= 0 && e <= 65535) this._bufferBuilder.append(205), this.pack_uint16(e);
                    else if (e >= -32768 && e <= 32767) this._bufferBuilder.append(209), this.pack_int16(e);
                    else if (e >= 0 && e <= 4294967295) this._bufferBuilder.append(206), this.pack_uint32(e);
                    else if (e >= -2147483648 && e <= 2147483647) this._bufferBuilder.append(210), this.pack_int32(e);
                    else if (e >= -9223372036854776e3 && e <= 9223372036854776e3) this._bufferBuilder.append(211), this.pack_int64(e);
                    else if (e >= 0 && e <= 18446744073709552e3) this._bufferBuilder.append(207), this.pack_uint64(e);
                    else throw new Error("Invalid integer");
                }
                pack_double(e) {
                    var t = 0;
                    _ = e < 0 && (t = 1, e = -e);
                    var f = Math.floor(Math.log(e) / Math.LN2);
                    var i = e / 2 ** f - 1;
                    var r = Math.floor(i * 2 ** 52);
                    var n = 2 ** 32;
                    var u = t << 31 | f + 1023 << 20 | r / n & 1048575;
                    var h = r % n;
                    this._bufferBuilder.append(203); this.pack_int32(u); this.pack_int32(h);
                }
                pack_object(e) {
                    var t = Object.keys(e);
                    var f = t.length;
                    if (f <= 15) this.pack_uint8(128 + f);
                    else if (f <= 65535) this._bufferBuilder.append(222), this.pack_uint16(f);
                    else if (f <= 4294967295) this._bufferBuilder.append(223), this.pack_uint32(f);
                    else throw new Error("Invalid length");
                    var i = r => {
                        if (r < t.length) {
                            let n = t[r];
                            if (e.hasOwnProperty(n)) {
                                this.pack(n);
                                let u = this.pack(e[n]);
                                if (u instanceof Promise) return u.then(() => i(r + 1));
                            }
                            return i(r + 1);
                        }
                    };
                    return i(0);
                }
                pack_uint8(e) {
                    this._bufferBuilder.append(e);
                }
                pack_uint16(e) {
                    this._bufferBuilder.append(e >> 8); this._bufferBuilder.append(e & 255);
                }
                pack_uint32(e) {
                    var t = e & 4294967295;
                    this._bufferBuilder.append((t & 4278190080) >>> 24); this._bufferBuilder.append((t & 16711680) >>> 16); this._bufferBuilder.append((t & 65280) >>> 8); this._bufferBuilder.append(t & 255);
                }
                pack_uint64(e) {
                    var t = e / 4294967296;
                    var f = e % 2 ** 32;
                    this._bufferBuilder.append((t & 4278190080) >>> 24); this._bufferBuilder.append((t & 16711680) >>> 16); this._bufferBuilder.append((t & 65280) >>> 8); this._bufferBuilder.append(t & 255); this._bufferBuilder.append((f & 4278190080) >>> 24); this._bufferBuilder.append((f & 16711680) >>> 16); this._bufferBuilder.append((f & 65280) >>> 8); this._bufferBuilder.append(f & 255);
                }
                pack_int8(e) {
                    this._bufferBuilder.append(e & 255);
                }
                pack_int16(e) {
                    this._bufferBuilder.append((e & 65280) >> 8); this._bufferBuilder.append(e & 255);
                }
                pack_int32(e) {
                    this._bufferBuilder.append(e >>> 24 & 255); this._bufferBuilder.append((e & 16711680) >>> 16); this._bufferBuilder.append((e & 65280) >>> 8); this._bufferBuilder.append(e & 255);
                }
                pack_int64(e) {
                    var t = Math.floor(e / 4294967296);
                    var f = e % 2 ** 32;
                    this._bufferBuilder.append((t & 4278190080) >>> 24); this._bufferBuilder.append((t & 16711680) >>> 16); this._bufferBuilder.append((t & 65280) >>> 8); this._bufferBuilder.append(t & 255); this._bufferBuilder.append((f & 4278190080) >>> 24); this._bufferBuilder.append((f & 16711680) >>> 16); this._bufferBuilder.append((f & 65280) >>> 8); this._bufferBuilder.append(f & 255);
                }
                constructor() {
                    this._bufferBuilder = new a();
                }
            };
    `)
	if err != nil {
		panic(err)
	}
	_, err = vm.RunString(`
    var aa = async (input) => {
        let e = new p; await e.pack(input);
        return new Uint8Array(e._bufferBuilder._parts);
    }
    `)
	if err != nil {
		panic(err)
	}
	processFunc, ok := goja.AssertFunction(vm.Get("aa"))
	if !ok {
		panic("processData is not a function")
	}
	result, err := processFunc(goja.Undefined(), vm.ToValue(str))
	if err != nil {
		panic(err)
	}
	if promise, ok := result.Export().(*goja.Promise); ok {
		switch promise.State() {
		case goja.PromiseStateRejected:
			panic("Promise rejected: " + promise.Result().String())
		case goja.PromiseStateFulfilled:
			return promise.Result().Export().([]byte)
		}
	}
	panic("Unexpected result type")
}
