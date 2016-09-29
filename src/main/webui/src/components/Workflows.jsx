'use strict';

//this is a component for constructing and running a workflow
//a workflow has the form
//
//+----------+        +----------+       +----------+                +-----------+     +----------+
//|          |        |          |       |          |                |           |     |          |
//|          |        |          |       |          |                |           |     |          |
//|  Volume 1+------->+Service 1 +------>+ Volume 2 +--->-+ + + + -->+ Service N +---->+ Volume N |
//|          |        |          |       |          |                |           |     |          |
//|          |        |          |       |          |                |           |     |          |
//|          |        |          |       |          |                |           |     |          |
//+----------+        +----------+       +----------+                +-----------+     +----------+
//
//
//each service is a job in the system, the user should be able to download the volumes in the chain
//The user can either upload their own data to Volume 1, or pick a a service which doesn't require any input(such as service
// fetching data from b2drop, these services are provided officially by GEF)

import React, {PropTypes} from 'react';
import bows from 'bows';
import _ from 'lodash';

const Workflow = () => (
    <div>Not implemented yet</div>
);

export default Workflow;
