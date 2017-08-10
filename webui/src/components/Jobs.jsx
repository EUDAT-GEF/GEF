import React from 'react';
import PropTypes from 'prop-types';
import {Row, Col} from 'react-bootstrap';
import {BootstrapTable, TableHeaderColumn} from 'react-bootstrap-table';
import {LinkContainer} from 'react-router-bootstrap'
import Job from './Job'

// const JobRow = ({job, title}) => (
//     <LinkContainer to={`/jobs/${job.ID}`}>
//         <Row>
//             <Col xs={12} sm={3} md={3}>{title}</Col>
//             <Col xs={12} sm={9} md={9}>{job.State.Status}</Col>
//         </Row>
//     </LinkContainer>
// );

// const Header = () => (
//     <div className="row table-head">
//         <div className="col-xs-12 col-sm-3">Job</div>
//         <div className="col-xs-12 col-sm-9">Status</div>
//     </div>
// );

const selectRowProp = {
    mode: 'checkbox'
};

var products = [{
    id: 1,
    name: "Product1",
    price: 120
}, {
    id: 2,
    name: "Product2",
    price: 80
}];

let order = 'desc';
let allJobs = [];
class Jobs extends React.Component {
// class MultiSelectTable extends React.Component {
    constructor(props) {
        super(props);

        this.options = {
            defaultSortName: 'name',  // default sort column name
            defaultSortOrder: 'desc'  // default sort order
        };


    }

    componentDidMount() {
        this.props.fetchJobs();
        this.props.fetchServices();


    }






    render() {

        allJobs = [];
        console.log(allJobs);
        console.log(allJobs.length);
        if (this.props.jobs) {
            return (
                <div>
                    <h3>Browse Jobs</h3>
                    <h4>All jobs</h4>
                    {this.props.jobs.map((job) => {
                        let service = null;
                        for (var i = 0; i < this.props.services.length; ++i) {
                            if (job.ServiceID == this.props.services[i].ID) {
                                service = this.props.services[i];
                                break;
                            }
                        }
                        let serviceName = (service && service.Name && service.Name.length) ? service.Name :
                            (service && service.ID && service.ID.length) ? service.ID : "unknown service";
                        let title = "Job from " + serviceName;


                        allJobs.push({"title": title, "id": job.ID, "name": job.ServiceID, "status": job.State.Status})
                        // let service = null;
                        // for (var i = 0; i < this.props.services.length; ++i) {
                        //     if (job.ServiceID == this.props.services[i].ID) {
                        //         service = this.props.services[i];
                        //         break;
                        //     }
                        // }
                        // const serviceName = (service && service.Name && service.Name.length) ? service.Name :
                        //     (service && service.ID && service.ID.length) ? service.ID : "unknown service";
                        // const title = "Job from " + serviceName;
                        //
                        // if (job.ID === this.props.params.id) {
                        //     return <Job key={job.ID} job={job} service={service} title={title}/>
                        // } else {
                        //     return <JobRow key={job.ID} job={job} title={title}/>
                        // }
                    })}

                    <div>
                        <BootstrapTable data={allJobs} selectRow={selectRowProp} options={this.options}>
                            <TableHeaderColumn width='20%' dataField='id' isKey dataSort>ID</TableHeaderColumn>
                            <TableHeaderColumn dataField='title' dataSort>Title</TableHeaderColumn>
                            <TableHeaderColumn width='20%' dataField='name' dataSort>Service ID</TableHeaderColumn>
                            <TableHeaderColumn width='20%' dataField='status'>Status</TableHeaderColumn>
                        </BootstrapTable>
                    </div>
                </div>
            );
        } else {
            return (
                <div><h4>No jobs found</h4></div>
            )
        }










    //
    //
    //     if (this.props.jobs) {
    //         return (
    //             <h3>Browse Jobs</h3>
    //         { this.props.jobs.map((job) => {
    //             console.log("TEXT");
    //             console.log(job);
    //             allJobs.push({"id": job.ID, "serviceID": job.ServiceID, "status": job.State.Status})
    //
    //
    //         })
    //         }
    //         // console.log(allJobs);
    //         // console.log(products);
    //
    //
    //         <div>
    //             <BootstrapTable data={products} selectRow={selectRowProp} options={this.options}>
    //                 <TableHeaderColumn dataField='id' isKey dataSort>Product ID</TableHeaderColumn>
    //                 <TableHeaderColumn dataField='name' dataSort>Product Name</TableHeaderColumn>
    //                 <TableHeaderColumn dataField='price'>Product Price</TableHeaderColumn>
    //             </BootstrapTable>
    //         </div>
    //     )
    //     }
    // }
    //
    //
    //




    }

}

Jobs.propTypes = {
//MultiSelectTable.propTypes = {
    jobs: PropTypes.array, // can be null
    fetchJobs: PropTypes.func.isRequired,
    services: PropTypes.array, // can be null
    fetchServices: PropTypes.func.isRequired,
    jobID: PropTypes.string
};

export default Jobs;
//export default MultiSelectTable;