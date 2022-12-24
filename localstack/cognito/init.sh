#!/bin/bash -e

user_pool_id=$(aws \
    --endpoint-url=http://localhost:9229 \
    cognito-idp create-user-pool --pool-name local-user-pool \
    --schema Name="tenant_id",AttributeDataType="String",DeveloperOnlyAttribute=false,Required=false,StringAttributeConstraints="{MinLength=1,MaxLength=256}" \
        Name="user_id",AttributeDataType="String",DeveloperOnlyAttribute=false,Required=false,StringAttributeConstraints="{MinLength=1,MaxLength=256}" \
        Name="user_role",AttributeDataType="String",DeveloperOnlyAttribute=false,Required=false,StringAttributeConstraints="{MinLength=1,MaxLength=256}" \
    --query UserPool.Id \
    | sed 's/"//g' \
    )

echo "user pool id is $user_pool_id"

client_id=$(aws \
    --endpoint-url=http://localhost:9229 \
    cognito-idp create-user-pool-client \
    --user-pool-id $user_pool_id \
    --client-name local-client \
    --query UserPoolClient.ClientId
    )

echo "client id is ${client_id}"
