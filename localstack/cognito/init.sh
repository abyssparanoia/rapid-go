#!/bin/bash -e

staff_user_pool_id=$(aws \
    --endpoint-url=http://localhost:9229 \
    cognito-idp create-user-pool --pool-name local-staff-user-pool \
    --schema Name="tenant_id",AttributeDataType="String",DeveloperOnlyAttribute=false,Required=false,StringAttributeConstraints="{MinLength=1,MaxLength=256}" \
        Name="staff_id",AttributeDataType="String",DeveloperOnlyAttribute=false,Required=false,StringAttributeConstraints="{MinLength=1,MaxLength=256}" \
        Name="staff_role",AttributeDataType="String",DeveloperOnlyAttribute=false,Required=false,StringAttributeConstraints="{MinLength=1,MaxLength=256}" \
    --query UserPool.Id \
    | sed 's/"//g' \
    )

echo "staff user pool id is $staff_user_pool_id"

staff_client_id=$(aws \
    --endpoint-url=http://localhost:9229 \
    cognito-idp create-user-pool-client \
    --user-pool-id $staff_user_pool_id \
    --client-name local-staff-client \
    --query UserPoolClient.ClientId
    )

echo "staff client id is ${staff_client_id}"

admin_user_pool_id=$(aws \
    --endpoint-url=http://localhost:9229 \
    cognito-idp create-user-pool --pool-name local-admin-user-pool \
    --schema Name="admin_id",AttributeDataType="String",DeveloperOnlyAttribute=false,Required=false,StringAttributeConstraints="{MinLength=1,MaxLength=256}" \
        Name="admin_role",AttributeDataType="String",DeveloperOnlyAttribute=false,Required=false,StringAttributeConstraints="{MinLength=1,MaxLength=256}" \
    --query UserPool.Id \
    | sed 's/"//g' \
    )

echo "admin user pool id is $admin_user_pool_id"

admin_client_id=$(aws \
    --endpoint-url=http://localhost:9229 \
    cognito-idp create-user-pool-client \
    --user-pool-id $admin_user_pool_id \
    --client-name local-admin-client \
    --query UserPoolClient.ClientId
    )

echo "admin client id is ${admin_client_id}"
