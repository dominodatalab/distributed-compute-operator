/*
 * Kubernetes
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: v1.18.2
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.dominodatalab.distributedcomputeoperator.dask.models;

import java.util.Objects;
import java.util.Arrays;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

/**
 * V1alpha1DaskClusterSpecSchedulerLifecyclePostStartExec
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:36:51.907Z[Etc/UTC]")
public class V1alpha1DaskClusterSpecSchedulerLifecyclePostStartExec {
  public static final String SERIALIZED_NAME_COMMAND = "command";
  @SerializedName(SERIALIZED_NAME_COMMAND)
  private List<String> command = null;


  public V1alpha1DaskClusterSpecSchedulerLifecyclePostStartExec command(List<String> command) {
    
    this.command = command;
    return this;
  }

  public V1alpha1DaskClusterSpecSchedulerLifecyclePostStartExec addCommandItem(String commandItem) {
    if (this.command == null) {
      this.command = new ArrayList<>();
    }
    this.command.add(commandItem);
    return this;
  }

   /**
   * Get command
   * @return command
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public List<String> getCommand() {
    return command;
  }


  public void setCommand(List<String> command) {
    this.command = command;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1DaskClusterSpecSchedulerLifecyclePostStartExec v1alpha1DaskClusterSpecSchedulerLifecyclePostStartExec = (V1alpha1DaskClusterSpecSchedulerLifecyclePostStartExec) o;
    return Objects.equals(this.command, v1alpha1DaskClusterSpecSchedulerLifecyclePostStartExec.command);
  }

  @Override
  public int hashCode() {
    return Objects.hash(command);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1DaskClusterSpecSchedulerLifecyclePostStartExec {\n");
    sb.append("    command: ").append(toIndentedString(command)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}
