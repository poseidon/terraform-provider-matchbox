# Group Resource

A Group matches (one or more) machines and declares a machine should be boot with a named `profile`.

```tf
resource "matchbox_group" "node1" {
  name = "node1"
  profile = "${matchbox_profile.myprofile.name}"
  selector = {
    mac = "52:54:00:a1:9c:ae"
  }
  metadata = {
    custom_variable = "machine_specific_value_here"
  }
}
```

## Argument Reference

* `name` - Unqiue name for the machine matcher
* `profile` - Name of a Matchbox profile
* `selector` - Map of hardware machine selectors. See [reserved selectors](https://matchbox.psdn.io/matchbox/#reserved-selectors). An empty selector becomes a global default group that matches machines.
* `metadata` - Map of group metadata (optional, seldom used)
