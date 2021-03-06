package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"strconv"
)

type SimpleChaincode struct {
}

type POLICY_ID_Holder struct {
	POLICY_IDs []string `json:"POLICY_IDs"`
}
type POLICY_ID_Holder_BeforeAppointment struct {
	POLICY_IDs []string `json:"POLICY_IDs"`
}
type HOSPITAL_NAME_Holder struct {
	HOSPITAL_NAMEs []string `json:"HOSPITAL_NAMEs"`
}
type PATIENT_HISTORY struct {
	PROPERTY_HISTORY_IDs []string `json:"PROPERTY_HISTORY_IDs"`
}

//Information to be stored about land in blockchain network
type Patient struct {
	PolicyId         string `json:"policyId"`
	City             string `json:"city"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Mobile           int64  `json:"mobile"`
	Hospital         string `json:"hospital"`
	AppointmentTime  string `json:"appointmentTime"`
	UnavailedBalance int64  `json:"unavailedBalance"`
	ClaimedAmount    int64  `json:"claimedAmount"`
}

type Policy struct {
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Address          string `json:"address"`
	Mobile           int64  `json:"mobile"`
	PolicyId         string `json:"policyId"`
	City             string `json:"city"`
	Pincode          string `json:"pincode"`
	UnavailedBalance int64  `json:"unavailedBalance"`
}

type User struct {
	Department string `json:"department"`
	UserName   string `json:"userName"`
	Password   string `json:"password"`
}

func main() {

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)

	}
}

//Init
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	var policyIdBytes []byte
	var policyId_beforeAppointment []byte
	var hospitalNameBytes []byte

	var policyIds POLICY_ID_Holder
	var hospitalNames HOSPITAL_NAME_Holder
	var policyIds_beforeAppointment POLICY_ID_Holder_BeforeAppointment

	policyIdBytes, err = json.Marshal(policyIds)
	hospitalNameBytes, err = json.Marshal(hospitalNames)
	policyId_beforeAppointment, err = json.Marshal(policyIds_beforeAppointment)

	if err != nil {
		return nil, errors.New("Error creating OWNER_ID_Holder record")
	}

	err = stub.PutState("policy_Ids", policyIdBytes)
	err = stub.PutState("hospital_Names", hospitalNameBytes)
	err = stub.PutState("policyIds_BA", policyId_beforeAppointment)
	return nil, nil

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "createAppointment" {

		return createAppointment(stub, args)
	}
	if function == "createPolicy" {

		return createPolicy(stub, args)
	}
	if function == "updateBalanceAPI" {

		return updateBalanceAPI(stub, args)
	}
	if function == "registerUser" {

		return registerUser(stub, args)
	}

	return nil, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "listAppointmentsByHospital" {

		return listAppointmentsByHospital(stub, args)
	}
	if function == "listAllAppointments" {

		return listAllAppointments(stub, args)
	}
	if function == "getPolicyDetails" {

		return getPolicyDetails(stub, args)
	}

	if function == "validateLogin" {

		return validateLogin(stub, args)
	}

	return nil, nil

}

func registerUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	user := User{}

	user.Department = args[0]
	user.UserName = args[1]
	user.Password = args[2]

	fmt.Println(user.Department, user.Password, user.UserName)

	userDetailsBytes, _ := stub.GetState(args[1])

	fmt.Println("userDetailsBytes = ", string(userDetailsBytes))
	userstr := string(userDetailsBytes)

	if len(userstr) != 0 {
		fmt.Println("Inside If")

		return []byte("User with username " + string(userDetailsBytes) + "already exists"), nil
	}
	fmt.Println("Inside else")

	userBytes, _ := json.Marshal(&user)
	fmt.Println("userBytes ", string(userBytes))
	err := stub.PutState(args[1], userBytes)

	if err != nil {
		return []byte("User registered Successfully"), nil
	}
	return nil, nil
}

func validateLogin(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	user := User{}

	userDetailsBytes, _ := stub.GetState(args[1])
	userStr := string(userDetailsBytes)

	if userStr == "" {
		return []byte("User with username " + user.UserName + "does not exists"), nil
	}

	json.Unmarshal(userDetailsBytes, &user)

	if user.Department == args[0] && user.UserName == args[1] && user.Password == args[2] {

		userDetailBytes, _ := json.Marshal(user)
		return userDetailBytes, nil

	}

	return nil, nil
}

func listAppointmentsByHospital(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	hospitalName := args[0]

	bytes, err := stub.GetState("policy_Ids")

	fmt.Println("Ids recieved", string(bytes))
	var policyIdHolder POLICY_ID_Holder
	err = json.Unmarshal(bytes, &policyIdHolder)

	result := []Patient{}

	var p Patient
	var resultBytes []byte

	for _, pat := range policyIdHolder.POLICY_IDs {

		fmt.Println("Inside for loop for getting PatientDetails for Hospital. Policy Id is  ", pat)

		p, err = retrievePatient(stub, pat)

		if p.Hospital == hospitalName {
			//temp, err = json.Marshal(p)

			if err == nil {

				result = append(result, p)
			}

		}

	}

	resultBytes, err = json.Marshal(result)

	if err != nil {
		return nil, err
	}

	return resultBytes, nil

}

func listAllAppointments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Inside list appointed Patients")
	var resultBytes []byte

	bytes, err := stub.GetState("policy_Ids")

	fmt.Println("Ids recieved", string(bytes))
	var policyIdHolder POLICY_ID_Holder
	err = json.Unmarshal(bytes, &policyIdHolder)

	result := []Patient{}

	var p Patient

	for _, pat := range policyIdHolder.POLICY_IDs {

		fmt.Println("Inside for loop for getting Property. Policy  Id is  ", pat)

		p, err = retrievePatient(stub, pat)

		if err == nil {
			result = append(result, p)
		}

	}
	resultBytes, err = json.Marshal(result)

	return resultBytes, nil
}

func retrievePatient(stub shim.ChaincodeStubInterface, patientId string) (Patient, error) {

	fmt.Println("Inside retrieve Patient")

	var p Patient

	bytes, err := stub.GetState(patientId)

	fmt.Println("Patient Id is ", patientId, "and Patient details are ", string(bytes))

	if err != nil {
		return p, errors.New("Patient not found")
	}

	err = json.Unmarshal(bytes, &p)

	return p, nil

}
func retrievePolicy(stub shim.ChaincodeStubInterface, policyId string) (Policy, error) {

	fmt.Println("Inside retrieve Patient")

	var p Policy

	bytes, err := stub.GetState(policyId)

	fmt.Println("Patient Id is ", policyId, "and Patient details are ", string(bytes))

	if err != nil {
		return p, errors.New("Patient not found")
	}

	err = json.Unmarshal(bytes, &p)

	return p, nil

}

func checkIfAppointmentExists(stub shim.ChaincodeStubInterface, args []string, methodName string) bool {

	result := false

	if methodName == "createApp" {

		policyId := args[0]

		bytes, _ := stub.GetState("policy_Ids")

		fmt.Println("Ids recieved", string(bytes))
		var policyIdHolder POLICY_ID_Holder
		json.Unmarshal(bytes, &policyIdHolder)

		for _, pat := range policyIdHolder.POLICY_IDs {

			fmt.Println("Inside for loop for getting Property. Policy  Id is  ", pat)
			if policyId == pat {
				result = true
			}

		}

	} else if methodName == "createPol" {

		policyId := args[4]

		bytes, _ := stub.GetState("policyIds_BA")
		fmt.Println("Ids recieved", string(bytes))
		var policyIdHolder POLICY_ID_Holder_BeforeAppointment
		json.Unmarshal(bytes, &policyIdHolder)

		for _, pat := range policyIdHolder.POLICY_IDs {

			fmt.Println("Inside for loop for getting Property. Policy  Id is  ", pat)

			if policyId == pat {
				result = true
			}

		}

	}

	return result
}

func updateBalanceAPI(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	policyId := args[0]

	claimedAmountInt, _ := strconv.ParseInt(args[1], 10, 64)

	patientDetailsBytes, err := stub.GetState(policyId)

	if err != nil {
		return nil, errors.New("error while getting Policy Id")
	}

	patientDetails := Patient{}

	err = json.Unmarshal(patientDetailsBytes, &patientDetails)

	unavailedBalance := patientDetails.UnavailedBalance

	if unavailedBalance >= claimedAmountInt {
		unavailedBalance = unavailedBalance - claimedAmountInt
	} else {
		unavailedBalance = 0
	}

	patientDetails.UnavailedBalance = unavailedBalance
	patientDetails.ClaimedAmount = claimedAmountInt

	patBytes, _ := json.Marshal(patientDetails)

	err = stub.PutState(policyId, patBytes)

	return patBytes, nil

}

func createAppointment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	appointmentsExists := checkIfAppointmentExists(stub, args, "createApp")
	fmt.Println("appointmentsExists : ", appointmentsExists)

	fmt.Println("Inside create appointment")

	patientDetails := Patient{}

	unavailedBalance, _ := strconv.ParseInt(args[7], 10, 64)
	claimedAmount, _ := strconv.ParseInt("0", 10, 64)
	mobile, _ := strconv.ParseInt(args[3], 10, 64)

	patientDetails.PolicyId = args[0]
	patientDetails.FirstName = args[1]
	patientDetails.LastName = args[2]
	patientDetails.Mobile = mobile
	patientDetails.Hospital = args[4]
	patientDetails.City = args[5]
	patientDetails.AppointmentTime = args[6]
	patientDetails.UnavailedBalance = unavailedBalance
	patientDetails.ClaimedAmount = claimedAmount

	patientDetailsBytes, err := json.Marshal(patientDetails)

	fmt.Println("Patient Details are : ", string(patientDetailsBytes))

	if err != nil {
		return nil, errors.New("Problem while saving Owner Details in BlockChain Network")

	}

	err = stub.PutState(patientDetails.PolicyId, patientDetailsBytes)

	//now owner has been added to block chain network, now we have to save the  Id as well

	bytes, err := stub.GetState("policy_Ids")

	var newPolicyId POLICY_ID_Holder

	err = json.Unmarshal(bytes, &newPolicyId)

	if err != nil {
		return nil, errors.New("error unmarshalling new Policy Id")
	}
	newPolicyId.POLICY_IDs = append(newPolicyId.POLICY_IDs, patientDetails.PolicyId)

	bytes, err = json.Marshal(newPolicyId)

	if err != nil {
		return nil, errors.New("error marshalling new Policy Id")
	}

	err = stub.PutState("policy_Ids", bytes)
	fmt.Println("Policy Id Saved is ", string(bytes))

	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	return nil, nil
}

func getPolicyDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	policyBytes, err := stub.GetState(args[0])

	if err != nil {
		return nil, errors.New("Error Unmarshalling")
	}

	return policyBytes, nil

}

func createPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	appointmentsExists := checkIfAppointmentExists(stub, args, "createPol")
	fmt.Println("appointmentsExists : ", appointmentsExists)

	if !appointmentsExists {
		fmt.Println("Inside create appointment")

		policyDetails := Policy{}

		unavailedBalance, _ := strconv.ParseInt(args[7], 10, 64)
		mobileNumber, _ := strconv.ParseInt(args[3], 10, 64)

		policyDetails.FirstName = args[0]
		policyDetails.LastName = args[1]
		policyDetails.Address = args[2]
		policyDetails.Mobile = mobileNumber
		policyDetails.PolicyId = args[4]
		policyDetails.City = args[5]
		policyDetails.Pincode = args[6]
		policyDetails.UnavailedBalance = unavailedBalance

		policyDetailsBytes, err := json.Marshal(policyDetails)

		fmt.Println("Policy Details are : ", string(policyDetailsBytes))

		if err != nil {
			return nil, errors.New("Problem while saving Owner Details in BlockChain Network")

		}

		err = stub.PutState(policyDetails.PolicyId, policyDetailsBytes)

		//now owner has been added to block chain network, now we have to save the  Id as well

		bytes, err := stub.GetState("policyIds_BA")

		var newPolicyId POLICY_ID_Holder_BeforeAppointment

		err = json.Unmarshal(bytes, &newPolicyId)

		if err != nil {
			return nil, errors.New("error unmarshalling new Policy Id")
		}
		newPolicyId.POLICY_IDs = append(newPolicyId.POLICY_IDs, policyDetails.PolicyId)

		bytes, err = json.Marshal(newPolicyId)

		if err != nil {
			return nil, errors.New("error marshalling new Policy Id")
		}

		err = stub.PutState("policyIds_BA", bytes)
		fmt.Println("Policy Id Saved is ", string(bytes))

		if err != nil {
			return nil, errors.New("Unable to put the state")
		}
	}
	return nil, nil
}
