mongosh <<EOF
use admin 
db.createUser( {user : "admin", pwd : "admin", roles : ["root"]})
exit
EOF

mongosh -u admin -p admin <<EOF
use admin 
db.createUser( {user : "app", pwd : "app", roles : ["root"]})
exit
EOF

mongosh -u admin -p admin <<EOF
use admin 
db.updateUser("app", {authenticationRestrictions: [{clientSource: ["127.0.0.1", "10.1.1.1"] }] })
exit
EOF


mongosh -u admin -p admin <<EOF
use admin 
db.system.users.updateOne ( {_id :"admin.app" }, { '\$unset': { authenticationRestrictions: 1 }})
exit
EOF

