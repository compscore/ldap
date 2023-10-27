# FTP

## Check Format

```yaml
- name:
  release:
    org: compscore
    repo: ftp
    tag: latest
  credentials:
    username:
    password:
  target:
  weight:
  options:
    ldaps:
```

## Parameters

| parameter  |          path           |   type   | default  | required | description                                     |
| :--------: | :---------------------: | :------: | :------: | :------: | :---------------------------------------------- |
|   `name`   |         `.name`         | `string` |   `""`   |  `true`  | `name of check (must be unique)`                |
|   `org`    |     `.release.org`      | `string` |   `""`   |  `true`  | `organization that check repository belongs to` |
|   `repo`   |     `.release.repo`     | `string` |   `""`   |  `true`  | `repository of the check`                       |
|   `tag`    |     `.release.tag`      | `string` | `latest` | `false`  | `tagged version of check`                       |
| `username` | `.credentials.username` | `string` |   `""`   | `false`  | `username of user (distinquished name format)`  |
| `password` | `.credentials.password` | `string` |   `""`   | `false`  | `default password of ldap user`                 |
|  `target`  |        `.target`        | `string` |   `""`   |  `true`  | `ldap server network location`                  |
|  `weight`  |        `.weight`        |  `int`   |   `0`    |  `true`  | `amount of points a successful check is worth`  |
|  `ldaps`   |    `.options.ldaps`     |  `bool`  | `false`  | `false`  | `use LDAPS instead of LDAP`                     |

## Examples

```yaml
- name: host_a-ldap
  release:
    org: compscore
    repo: ldap
    tag: latest
  credentials:
    username: cn=admin,dc=example,dc=local
    password: password
  target: 10.{{ .Team }}.1.1:389
  weight: 1
  options:
```

```yaml
- name: host_a-ldaps
  release:
    org: compscore
    repo: ldap
    tag: latest
  credentials:
    username: cn=admin,dc=example,dc=local
    password: password
  target: 10.{{ .Team }}.1.1:636
  weight: 1
  options:
    ldaps: true
```
