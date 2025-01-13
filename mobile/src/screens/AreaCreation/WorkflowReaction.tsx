import React, { useContext } from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { RootStackParamList } from '../../Navigation/navigate';
import { AppContext } from '../../context/AppContext';

type Props = NativeStackScreenProps<
  RootStackParamList,
  'WorkflowReactionScreen'
>;

const WorkflowReactionScreen = ({ navigation, route }: Props) => {
  const { actionId, actionOptions } = route.params;
  const { ipAddress, token } = useContext(AppContext);
  let service: any;
  const [name, setName] = React.useState('');

  const getService = async () => {
    try {
      const response = await fetch(
        `http://${ipAddress}:8080/api/v1/action/info/${actionId}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
        },
      );
      const serviceResponse = await fetch(
        `http://${ipAddress}:8080/api/v1/action/info/service/${actionId}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
        },
      );
      
      setName((await response.json())[0].name);
      service = (await serviceResponse.json())[0];
    } catch (error) {
      if (error.code === 401) {
        navigation.navigate('Login');
      }
      console.error('Error fetching service:', error);
    }
  };
  getService();

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Add Area</Text>
      <View style={styles.actionBox}>
        <Text style={styles.boxText}>{name}</Text>
        <TouchableOpacity
          style={[
            styles.addButtonDisabled,
            { backgroundColor: service?.color || '#ccc' },
          ]}>
          <Text style={styles.addTextDisabled}>Add</Text>
        </TouchableOpacity>
      </View>
      <View style={styles.line} />
      <View style={styles.actionBox}>
        <Text style={styles.boxText}>Reaction</Text>
        <TouchableOpacity
          style={styles.addButton}
          onPress={() =>
            navigation.navigate('AddReactionScreen', {
              actionId,
              actionOptions,
            })
          }>
          <Text style={styles.addTextDisabled}>Add</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: '#fff',
  },
  title: {
    fontSize: 24,
    marginBottom: 20,
  },
  actionBox: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: 'black',
    padding: 15,
    borderRadius: 8,
    width: '80%',
    marginBottom: 10,
  },
  reactionBox: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: 'gray',
    padding: 15,
    borderRadius: 8,
    width: '80%',
  },
  boxText: {
    color: '#fff',
    fontSize: 18,
  },
  addButton: {
    backgroundColor: 'white',
    paddingVertical: 5,
    paddingHorizontal: 15,
    borderRadius: 5,
  },
  addButtonDisabled: {
    paddingVertical: 5,
    paddingHorizontal: 15,
    borderRadius: 5,
  },
  addText: {
    color: 'black',
    fontSize: 16,
    fontWeight: 'bold',
  },
  addTextDisabled: {
    color: 'gray',
    fontSize: 16,
    fontWeight: 'bold',
  },
  line: {
    width: 2,
    height: 20,
    backgroundColor: 'black',
    marginVertical: 10,
  },
});

export default WorkflowReactionScreen;
