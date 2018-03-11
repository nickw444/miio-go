-- A Wireshark dissector written in for https://github.com/diacritic/wssdl.
-- Install by adding to ~/.config/wireshark/plugins along with wssdl.

local wssdl = require 'wssdl'

miio = wssdl.packet
{
  magic     : u16()
            : hex();
  length    : u16();
  unknown   : u32()
            : hex();
  deviceId  : u32()
            : hex();
  stamp     : u32();
  checksum  : bytes(16);
  data      : payload(magic);
}

wssdl.dissect {
  udp.port:set {
    [54321] = miio:proto('miio', 'Xiaomi Mi Home Binary Protocol')
  }
}
